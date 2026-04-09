package cmd_test

import (
	"path/filepath"
	"testing"
)

func TestPhase6CommonReadPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "common", "heartbeat.go"),
		filepath.Join("..", "internal", "domains", "common", "getGlobalConfig.go"),
		filepath.Join("..", "internal", "domains", "common", "getPrivacyPolicy.go"),
		filepath.Join("..", "internal", "domains", "common", "getTos.go"),
		filepath.Join("..", "internal", "domains", "common", "getStat.go"),
	})
}

func TestPhase6UserTicketAndSubscribeNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "user", "ticket", "createUserTicket.go"),
		filepath.Join("..", "internal", "domains", "user", "ticket", "createUserTicketFollow.go"),
		filepath.Join("..", "internal", "domains", "user", "ticket", "getUserTicketDetails.go"),
		filepath.Join("..", "internal", "domains", "user", "ticket", "getUserTicketList.go"),
		filepath.Join("..", "internal", "domains", "user", "ticket", "updateUserTicketStatus.go"),
		filepath.Join("..", "internal", "domains", "user", "subscribe", "querySubscribeGroupList.go"),
		filepath.Join("..", "internal", "domains", "user", "subscribe", "querySubscribeList.go"),
		filepath.Join("..", "internal", "domains", "user", "subscribe", "queryUserSubscribeNodeList.go"),
	})
}

func TestPhase6UserOrderNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "internal", "domains", "user", "order", "*.go")))
}

func TestPhase6UserPortalNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "internal", "domains", "user", "portal", "*.go")))
}

func TestPhase6UserReadPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "user", "user", "getDeviceList.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "getLoginLog.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "getOAuthMethods.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "getSubscribeLog.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "queryUserBalanceLog.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "queryUserCommissionLog.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "queryUserInfo.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "queryUserSubscribe.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "queryWithdrawalLog.go"),
	})
}

func TestPhase6UserProfileAndBindingWritesNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "user", "user", "updateUserNotify.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "updateUserPassword.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "updateUserRules.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "verifyEmail.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "updateBindEmail.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "updateBindMobile.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "unbindOAuth.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "bindTelegram.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "unbindTelegram.go"),
	})
}

func TestPhase6UserAffiliateAndDeviceFlowsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "user", "user", "queryUserAffiliate.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "queryUserAffiliateList.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "commissionWithdraw.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "unbindDevice.go"),
	})
}

func TestPhase6UserSubscriptionAndOAuthFlowsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "user", "user", "bindOAuth.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "bindOAuthCallback.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "calculateRemainingAmount.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "preUnsubscribe.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "resetUserSubscribeToken.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "updateUserSubscribeNote.go"),
		filepath.Join("..", "internal", "domains", "user", "user", "unsubscribe.go"),
	})
}

func TestPhase6AuthCheckPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "auth", "checkUser.go"),
		filepath.Join("..", "internal", "domains", "auth", "checkUserTelephone.go"),
	})
}

func TestPhase6AuthCorePathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "auth", "bindDevice.go"),
		filepath.Join("..", "internal", "domains", "auth", "deviceLogin.go"),
		filepath.Join("..", "internal", "domains", "auth", "resetPassword.go"),
		filepath.Join("..", "internal", "domains", "auth", "telephoneLogin.go"),
		filepath.Join("..", "internal", "domains", "auth", "telephoneResetPassword.go"),
		filepath.Join("..", "internal", "domains", "auth", "telephoneUserRegister.go"),
		filepath.Join("..", "internal", "domains", "auth", "userLogin.go"),
		filepath.Join("..", "internal", "domains", "auth", "userRegister.go"),
	})
}

func TestPhase6OAuthPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "auth", "oauth", "appleLoginCallback.go"),
		filepath.Join("..", "internal", "domains", "auth", "oauth", "oAuthLogin.go"),
		filepath.Join("..", "internal", "domains", "auth", "oauth", "oAuthLoginGetToken.go"),
	})
}

func TestPhase6SubscribeAndTelegramPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "subscribe", "subscribe.go"),
		filepath.Join("..", "internal", "domains", "telegram", "bot.go"),
		filepath.Join("..", "internal", "domains", "telegram", "telegram.go"),
	})
}

func TestPhase6AdminUserNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "internal", "domains", "admin", "user", "*.go")))
}
