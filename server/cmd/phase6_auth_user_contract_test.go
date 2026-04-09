package cmd_test

import (
	"path/filepath"
	"testing"
)

func TestPhase6CommonReadPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "common", "heartbeat.go"),
		filepath.Join("..", "services", "common", "getGlobalConfig.go"),
		filepath.Join("..", "services", "common", "getPrivacyPolicy.go"),
		filepath.Join("..", "services", "common", "getTos.go"),
		filepath.Join("..", "services", "common", "getStat.go"),
	})
}

func TestPhase6UserTicketAndSubscribeNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "user", "ticket", "createUserTicket.go"),
		filepath.Join("..", "services", "user", "ticket", "createUserTicketFollow.go"),
		filepath.Join("..", "services", "user", "ticket", "getUserTicketDetails.go"),
		filepath.Join("..", "services", "user", "ticket", "getUserTicketList.go"),
		filepath.Join("..", "services", "user", "ticket", "updateUserTicketStatus.go"),
		filepath.Join("..", "services", "user", "subscribe", "querySubscribeGroupList.go"),
		filepath.Join("..", "services", "user", "subscribe", "querySubscribeList.go"),
		filepath.Join("..", "services", "user", "subscribe", "queryUserSubscribeNodeList.go"),
	})
}

func TestPhase6UserOrderNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "user", "order", "*.go")))
}

func TestPhase6UserPortalNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "user", "portal", "*.go")))
}

func TestPhase6UserReadPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "user", "user", "getDeviceList.go"),
		filepath.Join("..", "services", "user", "user", "getLoginLog.go"),
		filepath.Join("..", "services", "user", "user", "getOAuthMethods.go"),
		filepath.Join("..", "services", "user", "user", "getSubscribeLog.go"),
		filepath.Join("..", "services", "user", "user", "queryUserBalanceLog.go"),
		filepath.Join("..", "services", "user", "user", "queryUserCommissionLog.go"),
		filepath.Join("..", "services", "user", "user", "queryUserInfo.go"),
		filepath.Join("..", "services", "user", "user", "queryUserSubscribe.go"),
		filepath.Join("..", "services", "user", "user", "queryWithdrawalLog.go"),
	})
}

func TestPhase6UserProfileAndBindingWritesNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "user", "user", "updateUserNotify.go"),
		filepath.Join("..", "services", "user", "user", "updateUserPassword.go"),
		filepath.Join("..", "services", "user", "user", "updateUserRules.go"),
		filepath.Join("..", "services", "user", "user", "verifyEmail.go"),
		filepath.Join("..", "services", "user", "user", "updateBindEmail.go"),
		filepath.Join("..", "services", "user", "user", "updateBindMobile.go"),
		filepath.Join("..", "services", "user", "user", "unbindOAuth.go"),
		filepath.Join("..", "services", "user", "user", "bindTelegram.go"),
		filepath.Join("..", "services", "user", "user", "unbindTelegram.go"),
	})
}

func TestPhase6UserAffiliateAndDeviceFlowsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "user", "user", "queryUserAffiliate.go"),
		filepath.Join("..", "services", "user", "user", "queryUserAffiliateList.go"),
		filepath.Join("..", "services", "user", "user", "commissionWithdraw.go"),
		filepath.Join("..", "services", "user", "user", "unbindDevice.go"),
	})
}

func TestPhase6UserSubscriptionAndOAuthFlowsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "user", "user", "bindOAuth.go"),
		filepath.Join("..", "services", "user", "user", "bindOAuthCallback.go"),
		filepath.Join("..", "services", "user", "user", "calculateRemainingAmount.go"),
		filepath.Join("..", "services", "user", "user", "preUnsubscribe.go"),
		filepath.Join("..", "services", "user", "user", "resetUserSubscribeToken.go"),
		filepath.Join("..", "services", "user", "user", "updateUserSubscribeNote.go"),
		filepath.Join("..", "services", "user", "user", "unsubscribe.go"),
	})
}

func TestPhase6AuthCheckPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "auth", "checkUser.go"),
		filepath.Join("..", "services", "auth", "checkUserTelephone.go"),
	})
}

func TestPhase6AuthCorePathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "auth", "bindDevice.go"),
		filepath.Join("..", "services", "auth", "deviceLogin.go"),
		filepath.Join("..", "services", "auth", "resetPassword.go"),
		filepath.Join("..", "services", "auth", "telephoneLogin.go"),
		filepath.Join("..", "services", "auth", "telephoneResetPassword.go"),
		filepath.Join("..", "services", "auth", "telephoneUserRegister.go"),
		filepath.Join("..", "services", "auth", "userLogin.go"),
		filepath.Join("..", "services", "auth", "userRegister.go"),
	})
}

func TestPhase6OAuthPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "auth", "oauth", "appleLoginCallback.go"),
		filepath.Join("..", "services", "auth", "oauth", "oAuthLogin.go"),
		filepath.Join("..", "services", "auth", "oauth", "oAuthLoginGetToken.go"),
	})
}

func TestPhase6SubscribeAndTelegramPathsNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "subscribe", "subscribe.go"),
		filepath.Join("..", "services", "telegram", "bot.go"),
		filepath.Join("..", "services", "telegram", "telegram.go"),
	})
}

func TestPhase6AdminUserNoLongerImportBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, globPhase6Targets(t, filepath.Join("..", "services", "admin", "user", "*.go")))
}
