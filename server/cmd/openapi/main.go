package main

import (
	"flag"
	"log"

	httpopenapi "github.com/perfect-panel/server/internal/platform/http/openapi"
)

func main() {
	outputDir := flag.String("o", "docs/openapi", "Output directory for spec files")
	flag.Parse()

	if err := httpopenapi.Export(*outputDir); err != nil {
		log.Fatal(err)
	}
}
