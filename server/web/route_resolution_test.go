package web

import (
	"testing"
	"testing/fstest"
)

func TestResolveEmbeddedRoutePrefersRealHTMLPage(t *testing.T) {
	staticFS := fstest.MapFS{
		"dashboard.html":       &fstest.MapFile{Data: []byte("<html>dashboard</html>")},
		"dashboard/index.html": &fstest.MapFile{Data: []byte("<html>dashboard index</html>")},
	}

	resolution := resolveEmbeddedRoute("/dashboard", staticFS)

	if resolution.kind != routeHTMLPage {
		t.Fatalf("expected HTML page, got %v", resolution.kind)
	}
	if resolution.filePath != "dashboard.html" {
		t.Fatalf("expected dashboard.html, got %q", resolution.filePath)
	}
}

func TestResolveEmbeddedRouteFallsBackToIndexWhenPageIsMissing(t *testing.T) {
	staticFS := fstest.MapFS{
		"favicon.ico": &fstest.MapFile{Data: []byte("ico")},
	}

	resolution := resolveEmbeddedRoute("/dashboard", staticFS)

	if resolution.kind != routeIndexFallback {
		t.Fatalf("expected index fallback, got %v", resolution.kind)
	}
}

func TestResolveEmbeddedRouteKeepsStaticAssets(t *testing.T) {
	staticFS := fstest.MapFS{
		"favicon.ico": &fstest.MapFile{Data: []byte("ico")},
	}

	resolution := resolveEmbeddedRoute("/favicon.ico", staticFS)

	if resolution.kind != routeStaticAsset {
		t.Fatalf("expected static asset, got %v", resolution.kind)
	}
	if resolution.filePath != "favicon.ico" {
		t.Fatalf("expected favicon.ico, got %q", resolution.filePath)
	}
}

func TestLegacyAdminRedirectTargetMapsLegacyPrefixToRuntimePath(t *testing.T) {
	target, ok := legacyAdminRedirectTarget("/admin/dashboard/servers", "/manage")
	if !ok {
		t.Fatal("expected legacy admin path to be redirected")
	}
	if target != "/manage/dashboard/servers" {
		t.Fatalf("expected /manage/dashboard/servers, got %q", target)
	}
}

func TestLegacyAdminRedirectTargetIgnoresCurrentAdminPath(t *testing.T) {
	if _, ok := legacyAdminRedirectTarget("/manage/dashboard", "/manage"); ok {
		t.Fatal("expected current admin path to skip legacy redirect")
	}
}
