package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	appweb "github.com/perfect-panel/server/web"
)

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	adminPath := envOrDefault("STATIC_SMOKE_ADMIN_PATH", "/manage")
	addr := envOrDefault("STATIC_SMOKE_ADDR", "127.0.0.1:4173")
	siteURL := envOrDefault("STATIC_SMOKE_SITE_URL", "http://"+addr)
	apiURL := envOrDefault("STATIC_SMOKE_API_URL", "")
	defaultLanguage := envOrDefault("STATIC_SMOKE_DEFAULT_LANGUAGE", "en-US")

	router := gin.New()
	router.Use(gin.Recovery())

	adminEnvVars := map[string]string{
		"VITE_ADMIN_PATH":         adminPath,
		"VITE_API_URL":            apiURL,
		"VITE_SITE_URL":           siteURL,
		"VITE_DEFAULT_LANGUAGE":   defaultLanguage,
		"VITE_DEFAULT_USER_EMAIL": "admin@ppanel.dev",
	}
	userEnvVars := map[string]string{
		"VITE_API_URL":          apiURL,
		"VITE_SITE_URL":         siteURL,
		"VITE_DEFAULT_LANGUAGE": defaultLanguage,
	}

	if err := appweb.RegisterStaticRoutes(router, adminPath, adminEnvVars, userEnvVars); err != nil {
		log.Fatalf("register static routes: %v", err)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("static smoke server listening on http://%s (admin path: %s)", addr, adminPath)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}
