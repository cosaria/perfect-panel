package initialize

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
)

func Site(ctx *svc.ServiceContext) {
	logger.Debug("initialize site config")
	configs, err := ctx.SystemModel.GetSiteConfig(context.Background())
	if err != nil {
		panic(err)
	}
	var siteConfig config.SiteConfig
	tool.SystemConfigSliceReflectToStruct(configs, &siteConfig)
	ctx.Config.Site = siteConfig
}
