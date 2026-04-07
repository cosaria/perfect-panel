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

// RegisterStaticRoutes sets up /admin/* SPA handler and user /* catch-all if embed is enabled.
// adminEnvVars and userEnvVars are injected into respective index.html as window.__ENV at startup.
func RegisterStaticRoutes(r *gin.Engine, adminEnvVars, userEnvVars map[string]string) error {
	if !embedEnabled {
		return nil
	}

	adminEnvJSON, err := json.Marshal(adminEnvVars)
	if err != nil {
		return fmt.Errorf("marshal admin env vars: %w", err)
	}

	userEnvJSON, err := json.Marshal(userEnvVars)
	if err != nil {
		return fmt.Errorf("marshal user env vars: %w", err)
	}

	// --- Admin frontend at /admin ---
	adminSub, err = fs.Sub(adminFS, "admin-dist")
	if err != nil {
		return fmt.Errorf("admin-dist sub fs: %w", err)
	}

	adminRaw, err := fs.ReadFile(adminFS, "admin-dist/index.html")
	if err != nil {
		return fmt.Errorf("read admin index.html: %w", err)
	}

	adminIndexHTML = bytes.Replace(adminRaw,
		[]byte("</head>"),
		[]byte(fmt.Sprintf("<script>window.__ENV=%s</script></head>", adminEnvJSON)),
		1,
	)

	adminFileServer := http.StripPrefix("/admin", http.FileServer(http.FS(adminSub)))

	r.GET("/admin", serveAdminIndex)
	r.GET("/admin/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		clean := path.Clean(strings.TrimPrefix(filepath, "/"))

		if !isStaticAsset(clean) {
			serveAdminIndex(c)
			return
		}

		f, err := adminSub.Open(clean)
		if err != nil {
			serveAdminIndex(c)
			return
		}
		defer f.Close()
		stat, err := f.Stat()
		if err != nil || stat.IsDir() {
			serveAdminIndex(c)
			return
		}

		if strings.HasPrefix(clean, "_next/static/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}

		adminFileServer.ServeHTTP(c.Writer, c.Request)
	})

	// --- User frontend at / (catch-all via NoRoute) ---
	userSub, err = fs.Sub(userFS, "user-dist")
	if err != nil {
		return fmt.Errorf("user-dist sub fs: %w", err)
	}

	userRaw, err := fs.ReadFile(userFS, "user-dist/index.html")
	if err != nil {
		return fmt.Errorf("read user index.html: %w", err)
	}

	userIndexHTML = bytes.Replace(userRaw,
		[]byte("</head>"),
		[]byte(fmt.Sprintf("<script>window.__ENV=%s</script></head>", userEnvJSON)),
		1,
	)

	userFileServer := http.FileServer(http.FS(userSub))

	r.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path

		// API paths should never fall through to the SPA — return proper JSON 404
		if strings.HasPrefix(reqPath, "/v1/") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
			return
		}

		clean := path.Clean(strings.TrimPrefix(reqPath, "/"))

		// SPA rule: only serve actual asset files directly.
		if clean != "" && isStaticAsset(clean) {
			f, err := userSub.Open(clean)
			if err == nil {
				defer f.Close()
				stat, err := f.Stat()
				if err == nil && !stat.IsDir() {
					if strings.HasPrefix(clean, "_next/static/") {
						c.Header("Cache-Control", "public, max-age=31536000, immutable")
					}
					userFileServer.ServeHTTP(c.Writer, c.Request)
					return
				}
			}
		}

		// All other requests get user index.html (SPA fallback)
		serveUserIndex(c)
	})

	return nil
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

func serveAdminIndex(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", adminIndexHTML)
}

func serveUserIndex(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", userIndexHTML)
}
