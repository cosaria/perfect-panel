package ads

import modelads "github.com/perfect-panel/server/models/ads"

// Deps holds the narrow admin/ads dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	AdsModel modelads.Model
}
