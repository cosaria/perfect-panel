package web

import (
	"bytes"
	"testing"
	"testing/fstest"
)

func TestResolveUserRoutePrefersExportedHTMLPage(t *testing.T) {
	staticFS := fstest.MapFS{
		"auth.html":       &fstest.MapFile{Data: []byte("<html><head></head><body>auth</body></html>")},
		"index.html":      &fstest.MapFile{Data: []byte("<html><head></head><body>home</body></html>")},
		"auth/index.html": &fstest.MapFile{Data: []byte("<html><head></head><body>auth nested</body></html>")},
	}

	got := resolveUserRoute("/auth", staticFS)
	if got.kind != routeHTMLPage {
		t.Fatalf("expected HTML page route, got %v", got.kind)
	}
	if got.filePath != "auth.html" {
		t.Fatalf("expected auth.html, got %q", got.filePath)
	}
}

func TestResolveUserRouteSupportsNestedIndexHTML(t *testing.T) {
	staticFS := fstest.MapFS{
		"bind/apple/index.html": &fstest.MapFile{Data: []byte("<html><head></head><body>bind apple</body></html>")},
	}

	got := resolveUserRoute("/bind/apple", staticFS)
	if got.kind != routeHTMLPage {
		t.Fatalf("expected HTML page route, got %v", got.kind)
	}
	if got.filePath != "bind/apple/index.html" {
		t.Fatalf("expected nested index.html, got %q", got.filePath)
	}
}

func TestResolveUserRouteKeepsUnknownPathAsIndexFallback(t *testing.T) {
	got := resolveUserRoute("/nonexistent-page", fstest.MapFS{})
	if got.kind != routeIndexFallback {
		t.Fatalf("expected index fallback, got %v", got.kind)
	}
}

func TestResolveUserRouteReturnsAPI404ForAPIPrefix(t *testing.T) {
	paths := []string{
		"/api/v1/common/site/config",
		"/v1/common/site/config",
	}

	for _, routePath := range paths {
		got := resolveUserRoute(routePath, fstest.MapFS{})
		if got.kind != routeAPI404 {
			t.Fatalf("path %q: expected API 404 route, got %v", routePath, got.kind)
		}
	}
}

func TestResolveUserRouteServesExistingStaticAsset(t *testing.T) {
	staticFS := fstest.MapFS{
		"assets/app-abc12345.js": &fstest.MapFile{Data: []byte("console.log('ok')")},
	}

	got := resolveUserRoute("/assets/app-abc12345.js", staticFS)
	if got.kind != routeStaticAsset {
		t.Fatalf("expected static asset route, got %v", got.kind)
	}
	if got.filePath != "assets/app-abc12345.js" {
		t.Fatalf("expected asset path to match, got %q", got.filePath)
	}
}

func TestInjectEnvIntoHTMLAddsWindowEnvScript(t *testing.T) {
	raw := []byte("<html><head><title>PPanel</title></head><body>auth</body></html>")
	envJSON := []byte(`{"VITE_API_URL":"http://localhost:8080"}`)

	got := injectEnvIntoHTML(raw, envJSON)
	if string(got) == string(raw) {
		t.Fatal("expected HTML to change after env injection")
	}
	want := `<script>window.__ENV={"VITE_API_URL":"http://localhost:8080"}</script></head>`
	if !bytes.Contains(got, []byte(want)) {
		t.Fatalf("expected injected env script %q, got %q", want, string(got))
	}
}

func TestInjectEnvIntoHTMLLeavesHeadlessDocumentUntouched(t *testing.T) {
	raw := []byte("<html><body>no head</body></html>")
	envJSON := []byte(`{"VITE_API_URL":"http://localhost:8080"}`)

	got := injectEnvIntoHTML(raw, envJSON)
	if string(got) != string(raw) {
		t.Fatalf("expected headless document to stay unchanged, got %q", string(got))
	}
}

func TestRewriteAdminHTMLBasePathUsesRuntimeAdminPath(t *testing.T) {
	raw := []byte(`<html><head></head><body><a href="/admin/dashboard">Dashboard</a><script src="/admin/assets/app-abc12345.js"></script></body></html>`)

	got := rewriteAdminHTMLBasePath(raw, "/manage")

	if bytes.Contains(got, []byte(`/admin/dashboard`)) {
		t.Fatalf("expected admin dashboard link to be rewritten, got %q", string(got))
	}
	if !bytes.Contains(got, []byte(`/manage/dashboard`)) {
		t.Fatalf("expected rewritten dashboard link, got %q", string(got))
	}
	if !bytes.Contains(got, []byte(`/manage/assets/app-abc12345.js`)) {
		t.Fatalf("expected rewritten vite asset path, got %q", string(got))
	}
}

func TestShouldUseImmutableAssetCache(t *testing.T) {
	testCases := []struct {
		path string
		want bool
	}{
		{path: "_next/static/chunks/app.js", want: false},
		{path: "assets/index-BqUBe2a_.js", want: true},
		{path: "assets/index-qFa1SYuY.css", want: true},
		{path: "favicon-BqUBe2a_.ico", want: true},
		{path: "images/logo.png", want: false},
		{path: "dashboard/index.html", want: false},
	}

	for _, tc := range testCases {
		if got := shouldUseImmutableAssetCache(tc.path); got != tc.want {
			t.Fatalf("path %q: expected %v, got %v", tc.path, tc.want, got)
		}
	}
}
