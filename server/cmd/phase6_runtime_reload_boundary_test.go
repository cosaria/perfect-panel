package cmd_test

import (
	"path/filepath"
	"testing"
)

func TestPhase6AuthMethodNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "admin", "authMethod", "getAuthMethodConfig.go"),
		filepath.Join("..", "services", "admin", "authMethod", "getAuthMethodList.go"),
		filepath.Join("..", "services", "admin", "authMethod", "getEmailPlatform.go"),
		filepath.Join("..", "services", "admin", "authMethod", "getSmsPlatform.go"),
		filepath.Join("..", "services", "admin", "authMethod", "testEmailSend.go"),
		filepath.Join("..", "services", "admin", "authMethod", "testSmsSend.go"),
		filepath.Join("..", "services", "admin", "authMethod", "updateAuthMethodConfig.go"),
	})
}

func TestPhase6SystemRuntimeSeamsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	targets := globPhase6Targets(t, filepath.Join("..", "services", "admin", "system", "*.go"))
	targets = append(targets, filepath.Join("..", "services", "admin", "tool", "restartSystem.go"))
	assertTargetsHaveNoBootstrapBoundaryDependency(t, targets)
}

func TestPhase6AdminLogNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "log", "*.go")))
}

func TestPhase6AdminSubscribeNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "subscribe", "*.go")))
}

func TestPhase6AdminServerNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "server", "*.go")))
}

func TestPhase6AdminPaymentNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "payment", "*.go")))
}

func TestPhase6AdminConsoleNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "console", "*.go")))
}

func TestPhase6AdminToolNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "tool", "*.go")))
}

func TestPhase6NotifyNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "platform", "http", "notify", "paymentNotify.go"),
		filepath.Join("..", "internal", "platform", "http", "notify", "ePayNotify.go"),
		filepath.Join("..", "internal", "platform", "http", "notify", "stripeNotify.go"),
		filepath.Join("..", "internal", "platform", "http", "notify", "alipayNotify.go"),
	})
}

func TestPhase6NodeNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "node", "getServerConfig.go"),
		filepath.Join("..", "services", "node", "getServerUserList.go"),
		filepath.Join("..", "services", "node", "pushOnlineUsers.go"),
		filepath.Join("..", "services", "node", "serverPushStatus.go"),
		filepath.Join("..", "services", "node", "serverPushUserTraffic.go"),
	})
}

func TestPhase6WorkerNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "jobs", "consumer_service.go"),
		filepath.Join("..", "internal", "jobs", "scheduler_service.go"),
		filepath.Join("..", "internal", "jobs", "email", "sendEmailLogic.go"),
		filepath.Join("..", "internal", "jobs", "email", "batchEmailLogic.go"),
		filepath.Join("..", "internal", "jobs", "sms", "sendSmsLogic.go"),
		filepath.Join("..", "internal", "jobs", "registry", "routes.go"),
		filepath.Join("..", "internal", "jobs", "order", "activateOrderLogic.go"),
		filepath.Join("..", "internal", "jobs", "order", "deferCloseOrderLogic.go"),
		filepath.Join("..", "internal", "jobs", "subscription", "checkSubscriptionLogic.go"),
		filepath.Join("..", "internal", "jobs", "task", "quotaLogic.go"),
		filepath.Join("..", "internal", "jobs", "task", "rateLogic.go"),
		filepath.Join("..", "internal", "jobs", "traffic", "resetTrafficLogic.go"),
		filepath.Join("..", "internal", "jobs", "traffic", "serverDataLogic.go"),
		filepath.Join("..", "internal", "jobs", "traffic", "trafficStatLogic.go"),
		filepath.Join("..", "internal", "jobs", "traffic", "trafficStatisticsLogic.go"),
	})
}

func TestPhase6InitializeNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "bootstrap", "configinit", "currency.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "device.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "email.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "init.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "invite.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "mobile.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "node.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "oauth.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "register.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "site.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "subscribe.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "telegram.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "verify.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "version.go"),
	})
}

func TestPhase6AdminMarketingNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "marketing", "*.go")))
}

func TestPhase6AdminApplicationNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "application", "*.go")))
}

func TestPhase6AdminTicketNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "ticket", "*.go")))
}

func TestPhase6AdminOrderNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "order", "*.go")))
}
