package seed

import (
	"encoding/json"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type systemSeed struct {
	Category string
	Key      string
	Value    string
	Type     string
	Desc     string
}

func Site(tx *gorm.DB) error {
	if tx == nil {
		return gorm.ErrInvalidDB
	}

	return tx.Transaction(func(tx *gorm.DB) error {
		if err := seedAuthMethods(tx); err != nil {
			return err
		}
		if err := seedSystemValues(tx); err != nil {
			return err
		}
		return nil
	})
}

func seedAuthMethods(tx *gorm.DB) error {
	rows := []auth.Auth{
		{Method: "email", Config: mustJSONString(auth.EmailAuthConfig{
			Platform:                   "smtp",
			PlatformConfig:             map[string]any{},
			EnableVerify:               false,
			EnableNotify:               false,
			EnableDomainSuffix:         false,
			DomainSuffixList:           "",
			VerifyEmailTemplate:        "",
			ExpirationEmailTemplate:    "",
			MaintenanceEmailTemplate:   "",
			TrafficExceedEmailTemplate: "",
		}), Enabled: boolPtr(true)},
		{Method: "mobile", Config: mustJSONString(auth.MobileAuthConfig{
			Platform:        "alibaba_cloud",
			PlatformConfig:  auth.AlibabaCloudConfig{},
			EnableWhitelist: false,
			Whitelist:       []string{},
		}), Enabled: boolPtr(false)},
		{Method: "apple", Config: `{"team_id":"","key_id":"","client_id":"","client_secret":"","redirect_url":""}`, Enabled: boolPtr(false)},
		{Method: "google", Config: `{"client_id":"","client_secret":"","redirect_url":""}`, Enabled: boolPtr(false)},
		{Method: "github", Config: `{"client_id":"","client_secret":"","redirect_url":""}`, Enabled: boolPtr(false)},
		{Method: "facebook", Config: `{"client_id":"","client_secret":"","redirect_url":""}`, Enabled: boolPtr(false)},
		{Method: "telegram", Config: mustJSONString(auth.TelegramAuthConfig{}), Enabled: boolPtr(false)},
		{Method: "device", Config: mustJSONString(auth.DeviceConfig{}), Enabled: boolPtr(false)},
	}

	for _, row := range rows {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedSystemValues(tx *gorm.DB) error {
	rows := []systemSeed{
		{Category: "site", Key: "SiteLogo", Value: "/favicon.svg", Type: "string", Desc: "Site Logo"},
		{Category: "site", Key: "SiteName", Value: "Perfect Panel", Type: "string", Desc: "Site Name"},
		{Category: "site", Key: "SiteDesc", Value: "PPanel is a pure, professional, and perfect open-source proxy panel tool, designed to be your ideal choice for learning and practical use.", Type: "string", Desc: "Site Description"},
		{Category: "site", Key: "Host", Value: "", Type: "string", Desc: "Site Host"},
		{Category: "site", Key: "Keywords", Value: "Perfect Panel,PPanel", Type: "string", Desc: "Site Keywords"},
		{Category: "site", Key: "CustomHTML", Value: "", Type: "string", Desc: "Custom HTML"},
		{Category: "tos", Key: "TosContent", Value: "Welcome to use Perfect Panel", Type: "string", Desc: "Terms of Service"},
		{Category: "tos", Key: "PrivacyPolicy", Value: "", Type: "string", Desc: "PrivacyPolicy"},
		{Category: "subscribe", Key: "SingleModel", Value: "false", Type: "bool", Desc: "是否单订阅模式"},
		{Category: "subscribe", Key: "SubscribePath", Value: "/api/subscribe", Type: "string", Desc: "订阅路径"},
		{Category: "subscribe", Key: "SubscribeDomain", Value: "", Type: "string", Desc: "订阅域名"},
		{Category: "subscribe", Key: "PanDomain", Value: "false", Type: "bool", Desc: "是否使用泛域名"},
		{Category: "verify", Key: "TurnstileSiteKey", Value: "", Type: "string", Desc: "TurnstileSiteKey"},
		{Category: "verify", Key: "TurnstileSecret", Value: "", Type: "string", Desc: "TurnstileSecret"},
		{Category: "verify", Key: "EnableLoginVerify", Value: "false", Type: "bool", Desc: "is enable login verify"},
		{Category: "verify", Key: "EnableRegisterVerify", Value: "false", Type: "bool", Desc: "is enable register verify"},
		{Category: "verify", Key: "EnableResetPasswordVerify", Value: "false", Type: "bool", Desc: "is enable reset password verify"},
		{Category: "server", Key: "NodeSecret", Value: "12345678", Type: "string", Desc: "node secret"},
		{Category: "server", Key: "NodePullInterval", Value: "10", Type: "int", Desc: "node pull interval"},
		{Category: "server", Key: "NodePushInterval", Value: "60", Type: "int", Desc: "node push interval"},
		{Category: "server", Key: "NodeMultiplierConfig", Value: "[]", Type: "string", Desc: "node multiplier config"},
		{Category: "invite", Key: "ForcedInvite", Value: "false", Type: "bool", Desc: "Forced invite"},
		{Category: "invite", Key: "ReferralPercentage", Value: "20", Type: "int", Desc: "Referral percentage"},
		{Category: "invite", Key: "OnlyFirstPurchase", Value: "false", Type: "bool", Desc: "Only first purchase"},
		{Category: "register", Key: "StopRegister", Value: "false", Type: "bool", Desc: "is stop register"},
		{Category: "register", Key: "EnableTrial", Value: "false", Type: "bool", Desc: "is enable trial"},
		{Category: "register", Key: "TrialSubscribe", Value: "", Type: "int", Desc: "Trial subscription"},
		{Category: "register", Key: "TrialTime", Value: "24", Type: "int", Desc: "Trial time"},
		{Category: "register", Key: "TrialTimeUnit", Value: "Hour", Type: "string", Desc: "Trial time unit"},
		{Category: "register", Key: "EnableIpRegisterLimit", Value: "false", Type: "bool", Desc: "is enable IP register limit"},
		{Category: "register", Key: "IpRegisterLimit", Value: "3", Type: "int", Desc: "IP Register Limit"},
		{Category: "register", Key: "IpRegisterLimitDuration", Value: "64", Type: "int", Desc: "IP Register Limit Duration (minutes)"},
		{Category: "currency", Key: "Currency", Value: "USD", Type: "string", Desc: "Currency"},
		{Category: "currency", Key: "CurrencySymbol", Value: "$", Type: "string", Desc: "Currency Symbol"},
		{Category: "currency", Key: "CurrencyUnit", Value: "USD", Type: "string", Desc: "Currency Unit"},
		{Category: "currency", Key: "AccessKey", Value: "", Type: "string", Desc: "Exchangerate Access Key"},
		{Category: "verify_code", Key: "VerifyCodeExpireTime", Value: "300", Type: "int", Desc: "Verify code expire time"},
		{Category: "verify_code", Key: "VerifyCodeLimit", Value: "15", Type: "int", Desc: "limits of verify code"},
		{Category: "verify_code", Key: "VerifyCodeInterval", Value: "60", Type: "int", Desc: "Interval of verify code"},
		{Category: "system", Key: "Version", Value: config.Version, Type: "string", Desc: "System Version"},
	}

	for _, row := range rows {
		record := system.System{
			Category: row.Category,
			Key:      row.Key,
			Value:    row.Value,
			Type:     row.Type,
			Desc:     row.Desc,
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

func mustJSONString(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func boolPtr(v bool) *bool {
	return &v
}
