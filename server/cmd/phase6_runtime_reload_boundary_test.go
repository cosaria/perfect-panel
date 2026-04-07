package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhase6AuthMethodNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "admin", "authMethod", "getAuthMethodConfig.go"),
		filepath.Join("..", "services", "admin", "authMethod", "getAuthMethodList.go"),
		filepath.Join("..", "services", "admin", "authMethod", "getEmailPlatform.go"),
		filepath.Join("..", "services", "admin", "authMethod", "getSmsPlatform.go"),
		filepath.Join("..", "services", "admin", "authMethod", "testEmailSend.go"),
		filepath.Join("..", "services", "admin", "authMethod", "testSmsSend.go"),
		filepath.Join("..", "services", "admin", "authMethod", "updateAuthMethodConfig.go"),
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6SystemRuntimeSeamsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "system", "*.go"))
	if err != nil {
		t.Fatalf("glob system targets: %v", err)
	}
	targets = append(targets, filepath.Join("..", "services", "admin", "tool", "restartSystem.go"))

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminLogNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "log", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/log targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminSubscribeNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "subscribe", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/subscribe targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminServerNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "server", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/server targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminPaymentNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "payment", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/payment targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminConsoleNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "console", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/console targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminToolNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "tool", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/tool targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6NotifyNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "notify", "paymentNotify.go"),
		filepath.Join("..", "services", "notify", "ePayNotify.go"),
		filepath.Join("..", "services", "notify", "stripeNotify.go"),
		filepath.Join("..", "services", "notify", "alipayNotify.go"),
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6NodeNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "node", "getServerConfig.go"),
		filepath.Join("..", "services", "node", "getServerUserList.go"),
		filepath.Join("..", "services", "node", "pushOnlineUsers.go"),
		filepath.Join("..", "services", "node", "serverPushStatus.go"),
		filepath.Join("..", "services", "node", "serverPushUserTraffic.go"),
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6WorkerNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "worker", "consumer_service.go"),
		filepath.Join("..", "worker", "scheduler_service.go"),
		filepath.Join("..", "worker", "email", "sendEmailLogic.go"),
		filepath.Join("..", "worker", "email", "batchEmailLogic.go"),
		filepath.Join("..", "worker", "sms", "sendSmsLogic.go"),
		filepath.Join("..", "worker", "registry", "routes.go"),
		filepath.Join("..", "worker", "order", "activateOrderLogic.go"),
		filepath.Join("..", "worker", "order", "deferCloseOrderLogic.go"),
		filepath.Join("..", "worker", "subscription", "checkSubscriptionLogic.go"),
		filepath.Join("..", "worker", "task", "quotaLogic.go"),
		filepath.Join("..", "worker", "task", "rateLogic.go"),
		filepath.Join("..", "worker", "traffic", "resetTrafficLogic.go"),
		filepath.Join("..", "worker", "traffic", "serverDataLogic.go"),
		filepath.Join("..", "worker", "traffic", "trafficStatLogic.go"),
		filepath.Join("..", "worker", "traffic", "trafficStatisticsLogic.go"),
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6InitializeNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "initialize", "currency.go"),
		filepath.Join("..", "initialize", "device.go"),
		filepath.Join("..", "initialize", "email.go"),
		filepath.Join("..", "initialize", "init.go"),
		filepath.Join("..", "initialize", "invite.go"),
		filepath.Join("..", "initialize", "mobile.go"),
		filepath.Join("..", "initialize", "node.go"),
		filepath.Join("..", "initialize", "oauth.go"),
		filepath.Join("..", "initialize", "register.go"),
		filepath.Join("..", "initialize", "site.go"),
		filepath.Join("..", "initialize", "subscribe.go"),
		filepath.Join("..", "initialize", "telegram.go"),
		filepath.Join("..", "initialize", "verify.go"),
		filepath.Join("..", "initialize", "version.go"),
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminMarketingNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "marketing", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/marketing targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminApplicationNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "application", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/application targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminTicketNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "ticket", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/ticket targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}

func TestPhase6AdminOrderNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "order", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/order targets: %v", err)
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}
