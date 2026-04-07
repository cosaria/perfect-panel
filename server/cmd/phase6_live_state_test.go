package cmd

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/models/node"
	modelsystem "github.com/perfect-panel/server/models/system"
	"github.com/perfect-panel/server/svc"
	"gorm.io/gorm"
)

func TestPhase6InitializeDepsPropagateToLiveState(t *testing.T) {
	t.Parallel()

	svcCtx := &svc.ServiceContext{}
	live := newLiveState(svcCtx)
	deps := newInitializeDeps(svcCtx, live)

	deps.SetExchangeRate(8.88)
	if got := live.ExchangeRate(); got != 8.88 {
		t.Fatalf("expected live exchange rate 8.88, got %v", got)
	}
	if got := svcCtx.ExchangeRate; got != 8.88 {
		t.Fatalf("expected service context exchange rate 8.88, got %v", got)
	}

	bot := &tgbotapi.BotAPI{}
	deps.SetTelegramBot(bot)
	if got := live.TelegramBot(); got != bot {
		t.Fatal("expected live telegram bot to update")
	}
	if got := svcCtx.TelegramBot; got != bot {
		t.Fatal("expected service context telegram bot to update")
	}

	manager := &node.Manager{}
	deps.SetNodeMultiplierManager(manager)
	if got := live.NodeMultiplierManager(); got != manager {
		t.Fatal("expected live node multiplier manager to update")
	}
	if got := svcCtx.NodeMultiplierManager; got != manager {
		t.Fatal("expected service context node multiplier manager to update")
	}
}

func TestPhase6WorkerDepsReadLiveState(t *testing.T) {
	t.Parallel()

	svcCtx := &svc.ServiceContext{}
	live := newLiveState(svcCtx)

	bot := &tgbotapi.BotAPI{}
	manager := &node.Manager{}
	live.SetTelegramBot(bot)
	live.SetNodeMultiplierManager(manager)

	orderDeps := newOrderWorkerDeps(svcCtx, live)
	if got := orderDeps.CurrentTelegramBot(); got != bot {
		t.Fatal("expected order worker deps to resolve telegram bot from live state")
	}

	trafficDeps := newTrafficWorkerDeps(svcCtx, live)
	if got := trafficDeps.CurrentNodeMultiplierManager(); got != manager {
		t.Fatal("expected traffic worker deps to resolve node multiplier manager from live state")
	}
}

func TestPhase6WorkerDepsLoadNodeMultiplierManagerFallback(t *testing.T) {
	t.Parallel()

	payload, err := json.Marshal([]node.TimePeriod{{
		StartTime:  "00:00.000",
		EndTime:    "23:59.000",
		Multiplier: 1.8,
	}})
	if err != nil {
		t.Fatalf("marshal multiplier config: %v", err)
	}

	svcCtx := &svc.ServiceContext{
		SystemModel: phase6SystemModelStub{
			findNodeMultiplierConfig: func(context.Context) (*modelsystem.System, error) {
				return &modelsystem.System{Value: string(payload)}, nil
			},
		},
	}
	live := newLiveState(svcCtx)

	trafficDeps := newTrafficWorkerDeps(svcCtx, live)
	manager, err := trafficDeps.ResolveNodeMultiplierManager(context.Background())
	if err != nil {
		t.Fatalf("expected loader fallback to succeed, got %v", err)
	}
	if manager == nil {
		t.Fatal("expected loader fallback to return a manager")
	}
	if got := manager.GetMultiplier(time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC)); got != 1.8 {
		t.Fatalf("expected loaded multiplier 1.8, got %v", got)
	}
	if live.NodeMultiplierManager() != manager {
		t.Fatal("expected loader fallback to update live state manager")
	}
	if svcCtx.NodeMultiplierManager != manager {
		t.Fatal("expected loader fallback to update service context manager")
	}
}

type phase6SystemModelStub struct {
	findNodeMultiplierConfig func(context.Context) (*modelsystem.System, error)
}

func (s phase6SystemModelStub) Insert(context.Context, *modelsystem.System) error { panic("unexpected Insert") }
func (s phase6SystemModelStub) FindOne(context.Context, int64) (*modelsystem.System, error) {
	panic("unexpected FindOne")
}
func (s phase6SystemModelStub) FindOneByKey(context.Context, string) (*modelsystem.System, error) {
	panic("unexpected FindOneByKey")
}
func (s phase6SystemModelStub) Update(context.Context, *modelsystem.System) error { panic("unexpected Update") }
func (s phase6SystemModelStub) Delete(context.Context, int64) error { panic("unexpected Delete") }
func (s phase6SystemModelStub) Transaction(context.Context, func(*gorm.DB) error) error {
	panic("unexpected Transaction")
}
func (s phase6SystemModelStub) GetSmsConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetSmsConfig")
}
func (s phase6SystemModelStub) GetSiteConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetSiteConfig")
}
func (s phase6SystemModelStub) GetSubscribeConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetSubscribeConfig")
}
func (s phase6SystemModelStub) GetRegisterConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetRegisterConfig")
}
func (s phase6SystemModelStub) GetVerifyConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetVerifyConfig")
}
func (s phase6SystemModelStub) GetNodeConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetNodeConfig")
}
func (s phase6SystemModelStub) GetInviteConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetInviteConfig")
}
func (s phase6SystemModelStub) GetTosConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetTosConfig")
}
func (s phase6SystemModelStub) GetCurrencyConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetCurrencyConfig")
}
func (s phase6SystemModelStub) GetVerifyCodeConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetVerifyCodeConfig")
}
func (s phase6SystemModelStub) GetLogConfig(context.Context) ([]*modelsystem.System, error) {
	panic("unexpected GetLogConfig")
}
func (s phase6SystemModelStub) UpdateNodeMultiplierConfig(context.Context, string) error {
	panic("unexpected UpdateNodeMultiplierConfig")
}
func (s phase6SystemModelStub) FindNodeMultiplierConfig(ctx context.Context) (*modelsystem.System, error) {
	if s.findNodeMultiplierConfig == nil {
		panic("unexpected FindNodeMultiplierConfig")
	}
	return s.findNodeMultiplierConfig(ctx)
}
