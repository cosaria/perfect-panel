package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhase6CommonReadPathsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "common", "heartbeat.go"),
		filepath.Join("..", "services", "common", "getGlobalConfig.go"),
		filepath.Join("..", "services", "common", "getPrivacyPolicy.go"),
		filepath.Join("..", "services", "common", "getTos.go"),
		filepath.Join("..", "services", "common", "getStat.go"),
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

func TestPhase6UserTicketAndSubscribeNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "user", "ticket", "createUserTicket.go"),
		filepath.Join("..", "services", "user", "ticket", "createUserTicketFollow.go"),
		filepath.Join("..", "services", "user", "ticket", "getUserTicketDetails.go"),
		filepath.Join("..", "services", "user", "ticket", "getUserTicketList.go"),
		filepath.Join("..", "services", "user", "ticket", "updateUserTicketStatus.go"),
		filepath.Join("..", "services", "user", "subscribe", "querySubscribeGroupList.go"),
		filepath.Join("..", "services", "user", "subscribe", "querySubscribeList.go"),
		filepath.Join("..", "services", "user", "subscribe", "queryUserSubscribeNodeList.go"),
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

func TestPhase6UserOrderNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "user", "order", "*.go"))
	if err != nil {
		t.Fatalf("glob user/order targets: %v", err)
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

func TestPhase6UserPortalNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "user", "portal", "*.go"))
	if err != nil {
		t.Fatalf("glob user/portal targets: %v", err)
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

func TestPhase6UserReadPathsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "user", "user", "getDeviceList.go"),
		filepath.Join("..", "services", "user", "user", "getLoginLog.go"),
		filepath.Join("..", "services", "user", "user", "getOAuthMethods.go"),
		filepath.Join("..", "services", "user", "user", "getSubscribeLog.go"),
		filepath.Join("..", "services", "user", "user", "queryUserBalanceLog.go"),
		filepath.Join("..", "services", "user", "user", "queryUserCommissionLog.go"),
		filepath.Join("..", "services", "user", "user", "queryUserInfo.go"),
		filepath.Join("..", "services", "user", "user", "queryUserSubscribe.go"),
		filepath.Join("..", "services", "user", "user", "queryWithdrawalLog.go"),
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

func TestPhase6UserProfileAndBindingWritesNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "user", "user", "updateUserNotify.go"),
		filepath.Join("..", "services", "user", "user", "updateUserPassword.go"),
		filepath.Join("..", "services", "user", "user", "updateUserRules.go"),
		filepath.Join("..", "services", "user", "user", "verifyEmail.go"),
		filepath.Join("..", "services", "user", "user", "updateBindEmail.go"),
		filepath.Join("..", "services", "user", "user", "updateBindMobile.go"),
		filepath.Join("..", "services", "user", "user", "unbindOAuth.go"),
		filepath.Join("..", "services", "user", "user", "bindTelegram.go"),
		filepath.Join("..", "services", "user", "user", "unbindTelegram.go"),
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

func TestPhase6UserAffiliateAndDeviceFlowsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "user", "user", "queryUserAffiliate.go"),
		filepath.Join("..", "services", "user", "user", "queryUserAffiliateList.go"),
		filepath.Join("..", "services", "user", "user", "commissionWithdraw.go"),
		filepath.Join("..", "services", "user", "user", "unbindDevice.go"),
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

func TestPhase6UserSubscriptionAndOAuthFlowsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "user", "user", "bindOAuth.go"),
		filepath.Join("..", "services", "user", "user", "bindOAuthCallback.go"),
		filepath.Join("..", "services", "user", "user", "calculateRemainingAmount.go"),
		filepath.Join("..", "services", "user", "user", "preUnsubscribe.go"),
		filepath.Join("..", "services", "user", "user", "resetUserSubscribeToken.go"),
		filepath.Join("..", "services", "user", "user", "updateUserSubscribeNote.go"),
		filepath.Join("..", "services", "user", "user", "unsubscribe.go"),
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

func TestPhase6AuthCheckPathsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "auth", "checkUser.go"),
		filepath.Join("..", "services", "auth", "checkUserTelephone.go"),
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

func TestPhase6AuthCorePathsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "auth", "bindDevice.go"),
		filepath.Join("..", "services", "auth", "deviceLogin.go"),
		filepath.Join("..", "services", "auth", "resetPassword.go"),
		filepath.Join("..", "services", "auth", "telephoneLogin.go"),
		filepath.Join("..", "services", "auth", "telephoneResetPassword.go"),
		filepath.Join("..", "services", "auth", "telephoneUserRegister.go"),
		filepath.Join("..", "services", "auth", "userLogin.go"),
		filepath.Join("..", "services", "auth", "userRegister.go"),
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

func TestPhase6OAuthPathsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "auth", "oauth", "appleLoginCallback.go"),
		filepath.Join("..", "services", "auth", "oauth", "oAuthLogin.go"),
		filepath.Join("..", "services", "auth", "oauth", "oAuthLoginGetToken.go"),
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

func TestPhase6SubscribeAndTelegramPathsNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "subscribe", "subscribe.go"),
		filepath.Join("..", "services", "telegram", "bot.go"),
		filepath.Join("..", "services", "telegram", "telegram.go"),
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

func TestPhase6AdminUserNoLongerImportServiceContext(t *testing.T) {
	t.Parallel()

	targets, err := filepath.Glob(filepath.Join("..", "services", "admin", "user", "*.go"))
	if err != nil {
		t.Fatalf("glob admin/user targets: %v", err)
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
