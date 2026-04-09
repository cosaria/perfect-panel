package handler

import "github.com/danielgtaylor/huma/v2"

var tagDescriptions = map[string]string{
	"ads":          "Admin advertising placement and promotion management.",
	"announcement": "Announcement and notice publishing APIs.",
	"application":  "Subscription application and template management APIs.",
	"auth":         "Authentication, registration, and credential recovery APIs.",
	"auth-method":  "Authentication provider and verification method configuration APIs.",
	"common":       "Shared unauthenticated runtime APIs used by multiple clients.",
	"console":      "Admin console and dashboard aggregate statistics APIs.",
	"coupon":       "Coupon creation, lookup, and lifecycle management APIs.",
	"document":     "Document and knowledge-base content APIs.",
	"log":          "Operational and audit log query APIs.",
	"marketing":    "Marketing automation and outbound campaign task APIs.",
	"oauth":        "OAuth provider bootstrap and callback token exchange APIs.",
	"order":        "Order placement, lookup, and lifecycle APIs.",
	"payment":      "Payment method discovery and payment configuration APIs.",
	"portal":       "Portal checkout and guest purchase workflow APIs.",
	"server":       "Node and server coordination APIs excluded from OpenAPI governance.",
	"subscribe":    "Subscription catalog and entitlement APIs.",
	"system":       "System configuration and runtime control APIs.",
	"ticket":       "Ticketing and support workflow APIs.",
	"tool":         "Operational utility and admin tooling APIs.",
	"user":         "Authenticated user profile, device, and account APIs.",
}

func governedAPIConfig(title string, version string, serverURL string, tags ...string) huma.Config {
	cfg := apiConfig(title, version)
	cfg.Servers = []*huma.Server{{URL: serverURL}}
	cfg.Components.SecuritySchemes = securitySchemes()
	cfg.Tags = governedOpenAPITags(tags...)
	return cfg
}

func governedOpenAPITags(tags ...string) []*huma.Tag {
	seen := map[string]struct{}{}
	out := make([]*huma.Tag, 0, len(tags))

	for _, tagName := range tags {
		if _, ok := seen[tagName]; ok {
			continue
		}
		seen[tagName] = struct{}{}

		tag := &huma.Tag{Name: tagName}
		if description, ok := tagDescriptions[tagName]; ok {
			tag.Description = description
		}
		out = append(out, tag)
	}

	return out
}

func governedSpecTags(tags ...string) []interface{} {
	items := make([]interface{}, 0, len(tags))
	for _, tag := range governedOpenAPITags(tags...) {
		items = append(items, map[string]interface{}{
			"name":        tag.Name,
			"description": tag.Description,
		})
	}
	return items
}
