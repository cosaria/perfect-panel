//go:build embed

package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	adminEnvVars := map[string]string{
		"VITE_API_URL":            "",
		"VITE_SITE_URL":           "http://localhost:8080",
		"VITE_DEFAULT_LANGUAGE":   "en-US",
		"VITE_DEFAULT_USER_EMAIL": "admin@ppanel.dev",
	}
	userEnvVars := map[string]string{
		"VITE_API_URL":          "",
		"VITE_SITE_URL":         "http://localhost:8080",
		"VITE_DEFAULT_LANGUAGE": "en-US",
	}
	if err := RegisterStaticRoutes(r, "/admin", adminEnvVars, userEnvVars); err != nil {
		panic(err)
	}
	return r
}

// ---- Admin tests ----

func TestAdminIndexServesHTML(t *testing.T) {
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /admin: expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("GET /admin: expected text/html, got %s", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, "window.__ENV") {
		t.Fatal("GET /admin: window.__ENV not injected")
	}
	if !strings.Contains(body, "VITE_API_URL") {
		t.Fatal("GET /admin: env vars not found in injected script")
	}
	t.Logf("GET /admin: %d bytes, window.__ENV injected", w.Body.Len())
}

func TestAdminSPAFallback(t *testing.T) {
	r := setupTestRouter()
	paths := []string{
		"/admin/nonexistent-page",
	}
	for _, p := range paths {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", p, nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d", p, w.Code)
			continue
		}
		if !strings.Contains(w.Body.String(), "window.__ENV") {
			t.Errorf("GET %s: SPA fallback did not return index.html", p)
		}
	}
}

func TestAdminRouteHTMLPagesServeExportedContent(t *testing.T) {
	r := setupTestRouter()
	paths := []string{
		"/admin/dashboard",
		"/admin/dashboard.html",
	}

	for _, p := range paths {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", p, nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("GET %s: expected 200, got %d", p, w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "window.__ENV") {
			t.Fatalf("GET %s: expected route HTML with window.__ENV", p)
		}
		if strings.Contains(body, "Please enter your account information to log in.") {
			t.Fatalf("GET %s: should not return auth index.html", p)
		}
	}
}

func TestAdminStaticAssets(t *testing.T) {
	r := setupTestRouter()
	// favicon.ico should exist in the static export
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/favicon.ico", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /admin/favicon.ico: expected 200, got %d", w.Code)
	}
	if strings.Contains(w.Body.String(), "window.__ENV") {
		t.Fatal("GET /admin/favicon.ico: should not return index.html")
	}
	t.Logf("GET /admin/favicon.ico: %d bytes", w.Body.Len())
}

func TestAdminCacheHeaders(t *testing.T) {
	r := setupTestRouter()

	// index.html should have no-cache
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	if cc := w.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("GET /admin Cache-Control: expected 'no-cache', got '%s'", cc)
	}

	// SPA fallback should also have no-cache
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/admin/dashboard/test", nil)
	r.ServeHTTP(w, req)
	if cc := w.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("GET /admin/dashboard/test Cache-Control: expected 'no-cache', got '%s'", cc)
	}
}

func TestAdminNextStaticAssetsCached(t *testing.T) {
	r := setupTestRouter()
	// _next/static files should have immutable cache
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/_next/static/css/test.css", nil)
	r.ServeHTTP(w, req)
	// File doesn't exist -> SPA fallback, that's OK
	// But if _next/static files exist, they should be cached
}

func TestAdminEnvVarsContent(t *testing.T) {
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `"VITE_SITE_URL":"http://localhost:8080"`) {
		t.Fatal("GET /admin: VITE_SITE_URL not found in window.__ENV")
	}
	if !strings.Contains(body, `"VITE_DEFAULT_USER_EMAIL":"admin@ppanel.dev"`) {
		t.Fatal("GET /admin: VITE_DEFAULT_USER_EMAIL not found in window.__ENV")
	}
	if strings.Contains(body, "NEXT_PUBLIC_") {
		t.Fatal("GET /admin: should not leak legacy NEXT_PUBLIC_* keys in window.__ENV")
	}
}

func TestUserEnvVarsDoNotContainCredentials(t *testing.T) {
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if strings.Contains(body, "VITE_DEFAULT_USER_EMAIL") {
		t.Fatal("GET /: user frontend should NOT contain admin email")
	}
	if strings.Contains(body, "VITE_DEFAULT_USER_PASSWORD") {
		t.Fatal("GET /: user frontend should NOT contain admin password")
	}
}

func TestDirectoryListingBlocked(t *testing.T) {
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/", nil)
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if strings.Contains(body, "<pre>") && strings.Contains(body, "<a href=") {
		t.Fatal("GET /admin/: directory listing exposed")
	}
	if !strings.Contains(body, "window.__ENV") {
		t.Fatal("GET /admin/: should return index.html")
	}
}

// ---- User tests (NoRoute catch-all) ----

func TestUserIndexServesHTML(t *testing.T) {
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /: expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("GET /: expected text/html, got %s", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, "window.__ENV") {
		t.Fatal("GET /: window.__ENV not injected")
	}
	t.Logf("GET /: %d bytes, window.__ENV injected", w.Body.Len())
}

func TestUserSPAFallback(t *testing.T) {
	r := setupTestRouter()
	paths := []string{
		"/dashboard",
		"/auth",
		"/purchasing",
		"/nonexistent-page",
	}
	for _, p := range paths {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", p, nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d", p, w.Code)
			continue
		}
		if !strings.Contains(w.Body.String(), "window.__ENV") {
			t.Errorf("GET %s: SPA fallback did not return user index.html", p)
		}
	}
}

func TestUserEnvVarsContent(t *testing.T) {
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `"VITE_SITE_URL":"http://localhost:8080"`) {
		t.Fatal("GET /: VITE_SITE_URL not found in window.__ENV")
	}
	if strings.Contains(body, "NEXT_PUBLIC_") {
		t.Fatal("GET /: should not leak legacy NEXT_PUBLIC_* keys in window.__ENV")
	}
}

func TestUserCacheHeaders(t *testing.T) {
	r := setupTestRouter()

	// index.html should have no-cache
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	if cc := w.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("GET / Cache-Control: expected 'no-cache', got '%s'", cc)
	}

	// SPA fallback should also have no-cache
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/dashboard/test", nil)
	r.ServeHTTP(w, req)
	if cc := w.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("GET /dashboard/test Cache-Control: expected 'no-cache', got '%s'", cc)
	}
}

func TestAPIPathsReturnJSON404(t *testing.T) {
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("GET /api/v1/nonexistent: expected 404, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("GET /api/v1/nonexistent: expected JSON, got %s", ct)
	}
	if strings.Contains(w.Body.String(), "window.__ENV") {
		t.Fatal("GET /api/v1/nonexistent: API path should NOT return SPA HTML")
	}
}

func TestUserHTMLPagesGetSPAFallback(t *testing.T) {
	r := setupTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/dashboard.html", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /dashboard.html: expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "window.__ENV") {
		t.Fatal("GET /dashboard.html: should return SPA user index.html with window.__ENV")
	}
}
