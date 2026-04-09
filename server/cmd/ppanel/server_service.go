package ppanel

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	appbootstrap "github.com/perfect-panel/server/internal/bootstrap/app"
	configinit "github.com/perfect-panel/server/internal/bootstrap/configinit"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/internal/domains/common/report"
	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/perfect-panel/server/modules/infra/proc"
	"github.com/perfect-panel/server/modules/infra/trace"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	handler "github.com/perfect-panel/server/internal/platform/http"
	"github.com/perfect-panel/server/internal/platform/http/middleware"
	"github.com/perfect-panel/server/web"
)

type Service struct {
	server     *http.Server
	svc        *appbootstrap.ServiceContext
	live       *appruntime.LiveState
	restarting atomic.Bool
}

func NewService(svc *appbootstrap.ServiceContext, live *appruntime.LiveState) *Service {
	if live == nil {
		live = newLiveState(svc)
	}
	return &Service{
		svc:  svc,
		live: live,
	}
}

func initServer(svc *appbootstrap.ServiceContext, live *appruntime.LiveState) *gin.Engine {

	// start init system config
	configinit.StartInitSystemConfig(newInitializeDeps(svc, live))
	// init gin server
	r := gin.Default()
	r.RemoteIPHeaders = []string{"X-Original-Forwarded-For", "X-Forwarded-For", "X-Real-IP"}
	// init session
	sessionStore, err := redis.NewStore(10, "tcp", svc.Config.Redis.Host, svc.Config.Redis.Pass, []byte(svc.Config.JwtAuth.AccessSecret))
	if err != nil {
		logger.Errorw("init session error", logger.Field("error", err.Error()))
		panic(err)
	}
	r.Use(sessions.Sessions("ppanel", sessionStore))
	// use cors middleware
	runtimeDeps := newRuntimeDeps(svc, live)
	r.Use(middleware.TraceMiddleware(), middleware.LoggerMiddleware(runtimeDeps), middleware.CorsMiddleware, gin.Recovery())

	// register handlers
	handler.RegisterHandlers(r, runtimeDeps)
	// register subscribe handler
	handler.RegisterSubscribeHandlers(r, runtimeDeps)
	// register telegram handler
	handler.RegisterTelegramHandlers(r, runtimeDeps)
	// register notify handler
	handler.RegisterNotifyHandlers(r, runtimeDeps)

	// register embedded frontends (only in production with -tags embed)
	adminPath := svc.Config.AdminPath
	if adminPath == "" {
		adminPath = "/admin"
	}

	adminEnvVars := map[string]string{
		"VITE_ADMIN_PATH":            adminPath,
		"VITE_SITE_URL":              svc.Config.Site.Host,
		"VITE_API_URL":               "", // same origin when embedded
		"VITE_DEFAULT_USER_EMAIL":    svc.Config.Administrator.Email,
		"VITE_DEFAULT_USER_PASSWORD": svc.Config.Administrator.Password,
		"VITE_DEFAULT_LANGUAGE":      "en-US",
	}
	userEnvVars := map[string]string{
		"VITE_SITE_URL":         svc.Config.Site.Host,
		"VITE_API_URL":          "", // same origin when embedded
		"VITE_DEFAULT_LANGUAGE": "en-US",
	}
	// debug mode: inject default credentials for user frontend
	if svc.Config.Debug {
		userEnvVars["VITE_DEFAULT_USER_EMAIL"] = svc.Config.Administrator.Email
		userEnvVars["VITE_DEFAULT_USER_PASSWORD"] = svc.Config.Administrator.Password
	}
	if err := web.RegisterStaticRoutes(r, adminPath, adminEnvVars, userEnvVars); err != nil {
		logger.Errorw("register static routes error", logger.Field("error", err.Error()))
	}

	return r
}

func (m *Service) Start() {
	m.run()
}

func (m *Service) run() {
	if m.svc == nil {
		panic("config file path is nil")
	}
	m.svc.Restart = m.Restart
	if m.live != nil {
		m.live.SetRestart(m.Restart)
	}

	// init service
	r := initServer(m.svc, m.live)
	// get server port
	port := m.svc.Config.Port
	host := m.svc.Config.Host
	// check gateway mode
	if report.IsGatewayMode() {
		// get free port
		freePort, err := report.ModulePort()
		if err != nil {
			logger.Errorf("get module port error: %s", err.Error())
			panic(err)
		}
		port = freePort
		host = "127.0.0.1"
		// register module
		err = report.RegisterModule(port)
		if err != nil {
			logger.Errorf("register module error: %s", err.Error())
			os.Exit(1)
		}
		logger.Infof("module registered on port %d", port)
	}

	serverAddr := fmt.Sprintf("%v:%d", host, port)
	m.server = &http.Server{
		Addr:    serverAddr,
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	trace.StartAgent(trace.Config{
		Name:    "ppanel",
		Sampler: 1.0,
		Batcher: "",
	})
	proc.AddShutdownListener(func() {
		trace.StopAgent()
	})
	logger.Infof("server start at %v", serverAddr)
	m.restarting.Store(false)
	if m.svc.Config.TLS.Enable {
		if err := m.server.ListenAndServeTLS(m.svc.Config.TLS.CertFile, m.svc.Config.TLS.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("server start error: %s", err.Error())
		}
	} else {
		if err := m.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("server start error: %s", err.Error())
		}
	}
}

func (m *Service) Stop() {
	if m.server == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.server.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown error: %s", err.Error())
	}
	logger.Info("server shutdown")
}

func (m *Service) Restart() error {
	if !m.restarting.CompareAndSwap(false, true) {
		return errors.New("server restart already in progress")
	}
	if m.server == nil {
		m.restarting.Store(false)
		return errors.New("server is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.server.Shutdown(ctx); err != nil {
		m.restarting.Store(false)
		logger.Errorf("server shutdown error: %v", err.Error())
		return err
	}
	m.server = nil
	logger.Info("server shutdown")
	go m.run()
	return nil
}
