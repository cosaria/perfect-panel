package nodeMultiplier

import (
	"testing"
	"time"
)

func TestNewNodeMultiplierManager(t *testing.T) {
	periods := []TimePeriod{
		{
			StartTime:  "23:00.000",
			EndTime:    "1:59.000",
			Multiplier: 1.2,
		},
		{
			StartTime:  "12:00.000",
			EndTime:    "13:59.000",
			Multiplier: 0.5,
		},
	}
	m := NewNodeMultiplierManager(periods)
	if len(m.Periods) != len(periods) {
		t.Fatalf("expected %d periods, got %d", len(periods), len(m.Periods))
	}

	if multiplier := m.GetMultiplier(time.Date(0, 1, 1, 0, 10, 0, 0, time.UTC)); multiplier != 1.2 {
		t.Fatalf("expected overnight multiplier 1.2, got %v", multiplier)
	}
	if multiplier := m.GetMultiplier(time.Date(0, 1, 1, 12, 30, 0, 0, time.UTC)); multiplier != 0.5 {
		t.Fatalf("expected midday multiplier 0.5, got %v", multiplier)
	}
	if multiplier := m.GetMultiplier(time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC)); multiplier != 1 {
		t.Fatalf("expected default multiplier 1, got %v", multiplier)
	}
}
