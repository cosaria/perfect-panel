package app

import (
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
)

func NewLogger(c config.Config) *logger.Logger {
	//log := logger.New(c.Logger)
	//// replace the default logger
	//logger = log
	//return log
	return nil
}
