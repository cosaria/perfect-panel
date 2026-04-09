package configinit

import (
	"time"

	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	"github.com/perfect-panel/server/internal/platform/persistence/schema/seed"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

func Migrate(deps Deps) {
	cfg := deps.currentConfig()
	now := time.Now()
	schemarevisions.RegisterEmbedded()
	if err := schema.Bootstrap(deps.DB, schema.SourceEmbedded); err != nil {
		logger.Errorf("[Migrate] schema bootstrap error: %v", err.Error())
		panic(err)
	}
	if err := seed.Site(deps.DB); err != nil {
		logger.Errorf("[Migrate] seed site error: %v", err.Error())
		panic(err)
	}
	if err := seed.Admin(deps.DB, cfg.Administrator.Email, cfg.Administrator.Password); err != nil {
		logger.Errorf("[Migrate] seed admin error: %v", err.Error())
		panic(err)
	}
	logger.Info("[Migrate] Database bootstrap complete, took " + time.Since(now).String())
}
