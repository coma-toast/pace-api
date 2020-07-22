package container

import (
	"sync"

	"github.com/coma-toast/pace-api/pkg/firestoredb"
	"github.com/coma-toast/pace-api/pkg/paceconfig"
	"github.com/coma-toast/pace-api/pkg/provider/user"
)

// Container exposes data providers
type Container interface {
	UserProvider() (user.Provider, error)
}

// Production is our production container for our external connections
type Production struct {
	config *paceconfig.Config
	// Providers
	userProvider *user.DatabaseProvider
	// Clients
	firestoreClient *firestoredb.Client
	// Mutex Locks
	userProviderMutex    *sync.Mutex
	firestoreClientMutex *sync.Mutex
}

// NewProduction builds a container with all of the config
func NewProduction(paceconfig *paceconfig.Config) Container {
	return &Production{
		config:               paceconfig,
		userProviderMutex:    &sync.Mutex{},
		firestoreClientMutex: &sync.Mutex{},
	}
}
