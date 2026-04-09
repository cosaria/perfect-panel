package cmd_test

import (
	"context"
	"errors"
	"testing"

	serverppanel "github.com/perfect-panel/server/cmd/ppanel"
	serverconfig "github.com/perfect-panel/server/config"
	appbootstrap "github.com/perfect-panel/server/internal/bootstrap/app"
	serverauth "github.com/perfect-panel/server/internal/domains/auth"
	servernodehandlers "github.com/perfect-panel/server/internal/domains/node"
	serversubscribe "github.com/perfect-panel/server/internal/domains/subscribe"
	servertelegram "github.com/perfect-panel/server/internal/domains/telegram"
	serveruserorder "github.com/perfect-panel/server/internal/domains/user/order"
	serverworker "github.com/perfect-panel/server/internal/jobs"
	servercache "github.com/perfect-panel/server/internal/platform/cache"
	serversyncx "github.com/perfect-panel/server/internal/platform/cache/syncx"
	serveraes "github.com/perfect-panel/server/internal/platform/crypto/aes"
	servernotify "github.com/perfect-panel/server/internal/platform/http/notify"
	serverresponse "github.com/perfect-panel/server/internal/platform/http/response"
	serveremail "github.com/perfect-panel/server/internal/platform/notify/email"
	servermigrate "github.com/perfect-panel/server/internal/platform/persistence/migrate"
	servernode "github.com/perfect-panel/server/internal/platform/persistence/node"
	serverauthmethod "github.com/perfect-panel/server/internal/platform/support/auth/authmethod"
	serverjwt "github.com/perfect-panel/server/internal/platform/support/auth/jwt"
	serverconf "github.com/perfect-panel/server/internal/platform/support/conf"
	servererrorx "github.com/perfect-panel/server/internal/platform/support/errorx"
	serverip "github.com/perfect-panel/server/internal/platform/support/ip"
	serverlimit "github.com/perfect-panel/server/internal/platform/support/limit"
	serverlogger "github.com/perfect-panel/server/internal/platform/support/logger"
	serverorm "github.com/perfect-panel/server/internal/platform/support/orm"
	serverproc "github.com/perfect-panel/server/internal/platform/support/proc"
	serverrescue "github.com/perfect-panel/server/internal/platform/support/rescue"
	serverrules "github.com/perfect-panel/server/internal/platform/support/rules"
	serverservice "github.com/perfect-panel/server/internal/platform/support/service"
	serverthreading "github.com/perfect-panel/server/internal/platform/support/threading"
	servertool "github.com/perfect-panel/server/internal/platform/support/tool"
	servertrace "github.com/perfect-panel/server/internal/platform/support/trace"
	serverturnstile "github.com/perfect-panel/server/internal/platform/support/verify/turnstile"
	serverxerr "github.com/perfect-panel/server/internal/platform/support/xerr"
)

func TestPhase1TopLevelPathsExist(t *testing.T) {
	var cfg serverconfig.Config
	var tempOrder serverconfig.TemporaryOrderInfo
	var multiplierPeriods []servernode.TimePeriod
	var ctx appbootstrap.ServiceContext
	var cache servercache.Cache
	var emailErr serveremail.ErrorInfo
	loadFn := serverconf.Load
	shutdownFn := serverproc.AddShutdownListener
	rescueFn := serverrescue.Recover
	limit := serversyncx.NewLimit(1)
	rule := serverrules.NewRule("DOMAIN,example.com,DIRECT", "fallback")
	periodLimit := serverlimit.NewPeriodLimit(1, 1, nil, "phase2:")
	token, err := serverjwt.NewJwtToken("secret", 1, 60, serverjwt.WithOption("sub", "user-1"))
	service := serverturnstile.New(serverturnstile.Config{Secret: "demo"})
	serviceGroup := serverservice.NewServiceGroup()
	routineGroup := serverthreading.NewRoutineGroup()
	multiplierManager := servernode.NewNodeMultiplierManager(multiplierPeriods)
	successPayload := serverresponse.Success(map[string]string{"phase": "2"})
	versionNumber := servertool.ExtractVersionNumber(serverconfig.Version)
	checkUserHandler := serverauth.CheckUserHandler(serverauth.Deps{})
	closeOrderHandler := serveruserorder.CloseOrderHandler(serveruserorder.Deps{})
	serverConfigHandler := servernodehandlers.GetServerConfigHandler(servernodehandlers.Deps{})
	paymentNotifyHandler := servernotify.PaymentNotifyHandler(servernotify.Deps{})
	subscribeHandler := serversubscribe.SubscribeHandler(serversubscribe.Deps{})
	telegramHandler := servertelegram.TelegramHandler(servertelegram.Deps{})
	consumerServiceCtor := serverworker.NewConsumerService
	schedulerServiceCtor := serverworker.NewSchedulerService
	executeFn := serverppanel.Execute
	errCode := serverxerr.NewErrMsg("boom")
	wrappedErr := servererrorx.Wrap(errors.New("inner"), "outer")
	logField := serverlogger.Field("phase", 2)
	parsedDSN := serverorm.ParseDSN("user:pass@tcp(localhost:3306)/ppanel")

	if servermigrate.NoChange == nil {
		t.Fatal("expected migrate package to expose NoChange")
	}

	if cfg.Port != 0 {
		t.Fatal("expected zero-value config in compile smoke test")
	}

	if serverconfig.Version == "" || serverconfig.BuildTime == "" || serverconfig.Repository == "" || serverconfig.ServiceName == "" {
		t.Fatal("expected config package to expose build metadata")
	}

	if serverconfig.CtxKeyUser == "" || serverconfig.LoginType == "" {
		t.Fatal("expected config package to expose request context keys")
	}

	if serverconfig.ParseVerifyType(uint8(serverconfig.Register)).String() != "register" {
		t.Fatal("expected config package to expose verify type helpers")
	}

	if tempOrder.OrderNo != "" || serverconfig.TempOrderCacheKey == "" {
		t.Fatal("expected config package to expose temporary order payload helpers")
	}

	if ctx.Redis != nil {
		t.Fatal("expected zero-value service context in compile smoke test")
	}

	if !limit.TryBorrow() {
		t.Fatal("expected syncx limit from internal/platform/cache to be usable")
	}

	if rule == nil {
		t.Fatal("expected rules package from internal/platform/support to parse a basic rule")
	}

	if token == "" || err != nil {
		t.Fatal("expected jwt package from internal/platform/support/auth to mint a token")
	}

	if serverauthmethod.Email == "" {
		t.Fatal("expected authmethod constants from internal/platform/support/auth")
	}

	if cache != nil {
		t.Fatal("expected zero-value cache interface in compile smoke test")
	}

	if _, _, err := serveraes.Encrypt([]byte("hello"), "secret"); err != nil {
		t.Fatal("expected aes package from internal/platform/crypto to encrypt test bytes")
	}

	if emailErr.Email != "" {
		t.Fatal("expected zero-value email worker struct")
	}

	if ips, err := serverip.GetIP("127.0.0.1"); err != nil || len(ips) != 1 {
		t.Fatal("expected util/ip package to resolve a direct IP")
	}

	if service == nil {
		t.Fatal("expected turnstile package from internal/platform/support/verify to create a service")
	}

	if loadFn == nil || shutdownFn == nil || rescueFn == nil {
		t.Fatal("expected infra function exports to be available")
	}

	if periodLimit == nil {
		t.Fatal("expected infra/limit package to construct a limiter")
	}

	if serviceGroup == nil || routineGroup == nil {
		t.Fatal("expected infra service and threading helpers to construct")
	}

	if multiplierManager == nil {
		t.Fatal("expected node model package to expose multiplier manager helpers")
	}

	if successPayload == nil || successPayload.Code != 200 {
		t.Fatal("expected internal/platform/http/response package to expose HTTP response helpers")
	}

	if versionNumber < 0 {
		t.Fatal("expected internal/platform/support/tool package to expose legacy utility helpers during phase 2")
	}

	if checkUserHandler == nil || closeOrderHandler == nil {
		t.Fatal("expected service packages to expose huma handler shims for phase 3 migration")
	}

	if serverConfigHandler == nil || paymentNotifyHandler == nil || subscribeHandler == nil || telegramHandler == nil {
		t.Fatal("expected non-huma entrypoints to be exposed from services packages during phase 3 migration")
	}

	if consumerServiceCtor == nil || schedulerServiceCtor == nil {
		t.Fatal("expected worker package to expose consumer and scheduler services during phase 4 migration")
	}

	if executeFn == nil {
		t.Fatal("expected command entry package to expose Execute from cmd/ppanel path")
	}

	if serverworker.SchedulerResetTraffic == "" || serverworker.ForthwithSendEmail == "" || serverworker.ForthwithActivateOrder == "" {
		t.Fatal("expected worker package to expose async task identifiers during phase 4 migration")
	}

	if errCode.GetErrMsg() != "boom" || wrappedErr == nil {
		t.Fatal("expected infra error helpers to be usable")
	}

	if logField.Key != "phase" || parsedDSN == nil {
		t.Fatal("expected infra logger and orm helpers to be usable")
	}

	if servertrace.TraceIDFromContext(context.Background()) != "" {
		t.Fatal("expected empty trace id from background context")
	}
}
