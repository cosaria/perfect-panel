package schema

import "fmt"

const (
	DefaultRevisionSource = "embedded"
	SourceEmbedded        = DefaultRevisionSource
	BaselineRevisionName  = "0001_baseline"
)

func NormalizeRevisionSource(source string) string {
	if source == "" {
		return DefaultRevisionSource
	}
	return source
}

func RevisionName(order int, slug string) string {
	return fmt.Sprintf("%04d_%s", order, slug)
}
