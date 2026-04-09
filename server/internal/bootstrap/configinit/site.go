package configinit

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/support/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/support/tool"
)

func Site(deps Deps) {
	logger.Debug("initialize site config")
	configs, err := deps.SystemModel.GetSiteConfig(context.Background())
	if err != nil {
		panic(err)
	}
	var siteConfig config.SiteConfig
	tool.SystemConfigSliceReflectToStruct(configs, &siteConfig)
	if deps.Config != nil {
		deps.Config.Site = siteConfig
	}
}
