package runtime

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/models/node"
)

// LiveState holds the mutable runtime values that can change after process start.
// Long-lived handlers and workers should read these through accessors instead of
// copying snapshots at construction time.
type LiveState struct {
	mu                    sync.RWMutex
	exchangeRate          ExchangeRateQuote
	restart               func() error
	telegramBot           *tgbotapi.BotAPI
	telegramPollerStop    func()
	nodeMultiplierManager *node.Manager
}

type ExchangeRateQuote struct {
	Version uint64
	From    string
	To      string
	Rate    float64
}

func NewLiveState() *LiveState {
	return &LiveState{}
}

func (s *LiveState) ExchangeRate() float64 {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.exchangeRate.Rate
}

func (s *LiveState) SetExchangeRate(rate float64) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.exchangeRate.Rate = rate
}

func (s *LiveState) ExchangeRateQuote() ExchangeRateQuote {
	if s == nil {
		return ExchangeRateQuote{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.exchangeRate
}

func (s *LiveState) PrepareExchangeRate(from, to string) uint64 {
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.exchangeRate.From != from || s.exchangeRate.To != to {
		s.exchangeRate.Version++
		s.exchangeRate.From = from
		s.exchangeRate.To = to
		s.exchangeRate.Rate = 0
	}
	return s.exchangeRate.Version
}

func (s *LiveState) StoreExchangeRate(version uint64, from, to string, rate float64) bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.exchangeRate.Version != version || s.exchangeRate.From != from || s.exchangeRate.To != to {
		return false
	}
	s.exchangeRate.Rate = rate
	return true
}

func (s *LiveState) Restart() func() error {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.restart
}

func (s *LiveState) SetRestart(restart func() error) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.restart = restart
}

func (s *LiveState) TelegramBot() *tgbotapi.BotAPI {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.telegramBot
}

func (s *LiveState) SetTelegramBot(bot *tgbotapi.BotAPI) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.telegramBot = bot
}

func (s *LiveState) NodeMultiplierManager() *node.Manager {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nodeMultiplierManager
}

func (s *LiveState) SetNodeMultiplierManager(manager *node.Manager) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodeMultiplierManager = manager
}

func (s *LiveState) SwapTelegramPoller(stop func()) (previous func()) {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	previous = s.telegramPollerStop
	s.telegramPollerStop = stop
	return previous
}
