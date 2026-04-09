package application

import (
	modelclient "github.com/perfect-panel/server/internal/platform/persistence/client"
	modelnode "github.com/perfect-panel/server/internal/platform/persistence/node"
)

// Deps holds the narrow admin application dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	ClientModel modelclient.Model
	NodeModel   modelnode.Model
}
