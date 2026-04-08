package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	adminIndexHTML []byte
	adminSub       fs.FS

	userIndexHTML []byte
	userSub       fs.FS
)

type routeKind int

const (
	routeAPI404 routeKind = iota
	routeStaticAsset
	routeHTMLPage
	routeIndexFallback
)

type routeResolution struct {
	kind     routeKind
	filePath string
}

// RegisterStaticRoutes sets up admin SPA handler at adminPath and user /* catch-all if embed is enabled.
// adminPath must start with "/" (e.g. "/admin", "/manage", "/secret-panel").
// adminEnvVars and userEnvVars are injected into respective index.html as window.__ENV at startup.
func RegisterStaticRoutes(r *gin.Engine, adminPath string, adminEnvVars, userEnvVars map[string]string) error {
	if !embedEnabled {
		return nil
	}

	// Normalize adminPath: ensure leading slash, no trailing slash
	adminPath = "/" + strings.Trim(adminPath, "/")

	adminEnvJSON, err := json.Marshal(adminEnvVars)
	if err != nil {
		return fmt.Errorf("marshal admin env vars: %w", err)
	}

	userEnvJSON, err := json.Marshal(userEnvVars)
	if err != nil {
		return fmt.Errorf("marshal user env vars: %w", err)
	}

	// --- Admin frontend at {adminPath} ---
	adminSub, err = fs.Sub(adminFS, "admin-dist")
	if err != nil {
		return fmt.Errorf("admin-dist sub fs: %w", err)
	}

	adminRaw, err := fs.ReadFile(adminFS, "admin-dist/index.html")
	if err != nil {
		return fmt.Errorf("read admin index.html: %w", err)
	}

	// Inject window.__ENV
	adminIndexHTML = injectEnvIntoHTML(rewriteAdminHTMLBasePath(adminRaw, adminPath), adminEnvJSON)

	adminFileServer := http.StripPrefix(adminPath, http.FileServer(http.FS(adminSub)))

	r.GET(adminPath, serveAdminIndex)
	r.GET(adminPath+"/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		resolution := resolveEmbeddedRoute(filepath, adminSub)
		switch resolution.kind {
		case routeStaticAsset:
			if strings.HasPrefix(resolution.filePath, "_next/static/") {
				c.Header("Cache-Control", "public, max-age=31536000, immutable")
			}
			adminFileServer.ServeHTTP(c.Writer, c.Request)
			return
		case routeHTMLPage:
			if serveAdminHTMLPage(c, adminSub, resolution.filePath, adminEnvJSON, adminPath) {
				return
			}
		}

		serveAdminIndex(c)
	})

	if adminPath != "/admin" {
		r.GET("/admin", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, adminPath)
		})
		r.GET("/admin/*filepath", func(c *gin.Context) {
			target, ok := legacyAdminRedirectTarget("/admin"+c.Param("filepath"), adminPath)
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
				return
			}
			if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
				target = fmt.Sprintf("%s?%s", target, rawQuery)
			}
			c.Redirect(http.StatusPermanentRedirect, target)
		})
	}

	// --- User frontend at / (catch-all via NoRoute) ---
	userSub, err = fs.Sub(userFS, "user-dist")
	if err != nil {
		return fmt.Errorf("user-dist sub fs: %w", err)
	}

	userRaw, err := fs.ReadFile(userFS, "user-dist/index.html")
	if err != nil {
		return fmt.Errorf("read user index.html: %w", err)
	}

	userIndexHTML = injectEnvIntoHTML(userRaw, userEnvJSON)

	userFileServer := http.FileServer(http.FS(userSub))

	r.NoRoute(func(c *gin.Context) {
		resolution := resolveUserRoute(c.Request.URL.Path, userSub)
		switch resolution.kind {
		case routeAPI404:
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
			return
		case routeStaticAsset:
			if strings.HasPrefix(resolution.filePath, "_next/static/") {
				c.Header("Cache-Control", "public, max-age=31536000, immutable")
			}
			userFileServer.ServeHTTP(c.Writer, c.Request)
			return
		case routeHTMLPage:
			if serveHTMLPage(c, userSub, resolution.filePath, userEnvJSON) {
				return
			}
			serveUserIndex(c)
			return
		}

		serveUserIndex(c)
	})

	return nil
}

func resolveUserRoute(reqPath string, staticFS fs.FS) routeResolution {
	if strings.HasPrefix(reqPath, "/api/") || strings.HasPrefix(reqPath, "/v1/") {
		return routeResolution{kind: routeAPI404}
	}

	return resolveEmbeddedRoute(reqPath, staticFS)
}

func resolveEmbeddedRoute(reqPath string, staticFS fs.FS) routeResolution {
	clean := path.Clean(strings.TrimPrefix(reqPath, "/"))
	if clean == "." {
		clean = ""
	}

	if clean != "" && isStaticAsset(clean) && fileExists(staticFS, clean) {
		return routeResolution{kind: routeStaticAsset, filePath: clean}
	}

	if clean != "" {
		if htmlPath, ok := resolveRouteHTMLPath(clean, staticFS); ok {
			return routeResolution{kind: routeHTMLPage, filePath: htmlPath}
		}
	}

	return routeResolution{kind: routeIndexFallback}
}

func resolveRouteHTMLPath(clean string, staticFS fs.FS) (string, bool) {
	candidates := make([]string, 0, 3)
	if path.Ext(clean) == ".html" {
		candidates = append(candidates, clean)
	}
	candidates = append(candidates, clean+".html", path.Join(clean, "index.html"))

	for _, candidate := range candidates {
		if fileExists(staticFS, candidate) {
			return candidate, true
		}
	}

	return "", false
}

func fileExists(staticFS fs.FS, name string) bool {
	f, err := staticFS.Open(name)
	if err != nil {
		return false
	}
	defer func() {
		_ = f.Close()
	}()
	stat, err := f.Stat()
	return err == nil && !stat.IsDir()
}

// isStaticAsset returns true if the path looks like a static asset
// (has a file extension that is NOT .html).
func isStaticAsset(p string) bool {
	ext := path.Ext(p)
	if ext == "" || ext == ".html" {
		return false
	}
	return true
}

func injectEnvIntoHTML(raw, envJSON []byte) []byte {
	return bytes.Replace(
		raw,
		[]byte("</head>"),
		[]byte(fmt.Sprintf("<script>window.__ENV=%s</script></head>", envJSON)),
		1,
	)
}

func rewriteAdminHTMLBasePath(raw []byte, adminPath string) []byte {
	if adminPath == "/admin" {
		return raw
	}

	rewritten := bytes.ReplaceAll(raw, []byte(`"/admin/`), []byte(`"`+adminPath+`/`))
	rewritten = bytes.ReplaceAll(rewritten, []byte(`"/admin"`), []byte(`"`+adminPath+`"`))
	rewritten = bytes.ReplaceAll(rewritten, []byte(`href="/admin/`), []byte(`href="`+adminPath+`/`))
	rewritten = bytes.ReplaceAll(rewritten, []byte(`href="/admin"`), []byte(`href="`+adminPath+`"`))
	rewritten = bytes.ReplaceAll(rewritten, []byte(`src="/admin/`), []byte(`src="`+adminPath+`/`))
	rewritten = bytes.ReplaceAll(rewritten, []byte(`src="/admin"`), []byte(`src="`+adminPath+`"`))
	return rewritten
}

func legacyAdminRedirectTarget(reqPath, adminPath string) (string, bool) {
	if adminPath == "" || adminPath == "/admin" {
		return "", false
	}
	if reqPath != "/admin" && !strings.HasPrefix(reqPath, "/admin/") {
		return "", false
	}

	suffix := strings.TrimPrefix(reqPath, "/admin")
	return adminPath + suffix, true
}

func serveHTMLPage(c *gin.Context, staticFS fs.FS, filePath string, envJSON []byte) bool {
	raw, err := fs.ReadFile(staticFS, filePath)
	if err != nil {
		return false
	}
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", injectEnvIntoHTML(raw, envJSON))
	return true
}

func serveAdminHTMLPage(c *gin.Context, staticFS fs.FS, filePath string, envJSON []byte, adminPath string) bool {
	raw, err := fs.ReadFile(staticFS, filePath)
	if err != nil {
		return false
	}
	c.Header("Cache-Control", "no-cache")
	c.Data(
		http.StatusOK,
		"text/html; charset=utf-8",
		injectEnvIntoHTML(rewriteAdminHTMLBasePath(raw, adminPath), envJSON),
	)
	return true
}

func serveAdminIndex(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", adminIndexHTML)
}

func serveUserIndex(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", userIndexHTML)
}
