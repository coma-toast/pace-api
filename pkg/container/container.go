package container

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/coma-toast/pace-api/pkg/paceconfig"
	"github.com/coma-toast/pace-api/pkg/provider/company"
	"github.com/coma-toast/pace-api/pkg/provider/contact"
	"github.com/coma-toast/pace-api/pkg/provider/firestoredb"
	"github.com/coma-toast/pace-api/pkg/provider/inspection"
	"github.com/coma-toast/pace-api/pkg/provider/inventory"
	"github.com/coma-toast/pace-api/pkg/provider/project"
	"github.com/coma-toast/pace-api/pkg/provider/user"
	"google.golang.org/api/option"
)

// Container exposes data providers
type Container interface {
	UserProvider() (user.Provider, error)
	ContactProvider() (contact.Provider, error)
	CompanyProvider() (company.Provider, error)
	ProjectProvider() (project.Provider, error)
	InspectionProvider() (inspection.Provider, error)
	InventoryProvider() (inventory.Provider, error)
}

// Production is our production container for our external connections
type Production struct {
	config *paceconfig.Config
	// Providers
	userProvider       *user.DatabaseProvider
	contactProvider    *contact.DatabaseProvider
	companyProvider    *company.DatabaseProvider
	projectProvider    *project.DatabaseProvider
	inspectionProvider *inspection.DatabaseProvider
	inventoryProvider  *inventory.DatabaseProvider
	// Clients
	firestoreClient *firestore.Client
	// Mutex Locks
	userProviderMutex       *sync.Mutex
	contactProviderMutex    *sync.Mutex
	companyProviderMutex    *sync.Mutex
	projectProviderMutex    *sync.Mutex
	inspectionProviderMutex *sync.Mutex
	inventoryProviderMutex  *sync.Mutex
	firestoreClientMutex    *sync.Mutex
}

// UserProvider provides the user provider
func (p Production) UserProvider() (user.Provider, error) {
	if p.userProvider != nil {
		return p.userProvider, nil
	}

	firestoreConnection, err := p.getFirestoreConnection()
	if err != nil {
		return nil, err
	}

	p.userProvider = &user.DatabaseProvider{
		SharedProvider: &firestoredb.DatabaseProvider{
			Database:   firestoreConnection,
			Collection: "users",
		},
	}

	return p.userProvider, nil
}

// ContactProvider provides the contact provider
func (p Production) ContactProvider() (contact.Provider, error) {
	if p.contactProvider != nil {
		return p.contactProvider, nil
	}

	firestoreConnection, err := p.getFirestoreConnection()
	if err != nil {
		return nil, err
	}
	p.contactProvider = &contact.DatabaseProvider{
		SharedProvider: &firestoredb.DatabaseProvider{
			Database:   firestoreConnection,
			Collection: "contacts",
		}}

	return p.contactProvider, nil
}

// CompanyProvider provides the Company provider
func (p Production) CompanyProvider() (company.Provider, error) {
	if p.companyProvider != nil {
		return p.companyProvider, nil
	}

	firestoreConnection, err := p.getFirestoreConnection()
	if err != nil {
		return nil, err
	}

	p.companyProvider = &company.DatabaseProvider{
		SharedProvider: &firestoredb.DatabaseProvider{
			Database:   firestoreConnection,
			Collection: "company",
		},
	}
	return p.companyProvider, nil
}

// ProjectProvider provides the Company provider
func (p Production) ProjectProvider() (project.Provider, error) {
	if p.projectProvider != nil {
		return p.projectProvider, nil
	}

	firestoreConnection, err := p.getFirestoreConnection()
	if err != nil {
		return nil, err
	}

	p.projectProvider = &project.DatabaseProvider{
		SharedProvider: &firestoredb.DatabaseProvider{
			Database:   firestoreConnection,
			Collection: "projects",
		},
	}

	return p.projectProvider, nil
}

// InspectionProvider provides the Company provider
func (p Production) InspectionProvider() (inspection.Provider, error) {
	if p.inspectionProvider != nil {
		return p.inspectionProvider, nil
	}

	firestoreConnection, err := p.getFirestoreConnection()
	if err != nil {
		return nil, err
	}

	p.inspectionProvider = &inspection.DatabaseProvider{
		SharedProvider: &firestoredb.DatabaseProvider{
			Database:   firestoreConnection,
			Collection: "inspections",
		},
	}

	return p.inspectionProvider, nil
}

// InventoryProvider provides the Company provider
func (p Production) InventoryProvider() (inventory.Provider, error) {
	if p.inventoryProvider != nil {
		return p.inventoryProvider, nil
	}

	firestoreConnection, err := p.getFirestoreConnection()
	if err != nil {
		return nil, err
	}

	p.inventoryProvider = &inventory.DatabaseProvider{
		SharedProvider: &firestoredb.DatabaseProvider{
			Database:   firestoreConnection,
			Collection: "inventory",
		},
	}

	return p.inventoryProvider, nil
}

// NewProduction builds a container with all of the config
func NewProduction(paceconfig *paceconfig.Config) Container {
	return &Production{
		config:                  paceconfig,
		userProviderMutex:       &sync.Mutex{},
		contactProviderMutex:    &sync.Mutex{},
		companyProviderMutex:    &sync.Mutex{},
		projectProviderMutex:    &sync.Mutex{},
		inspectionProviderMutex: &sync.Mutex{},
		inventoryProviderMutex:  &sync.Mutex{},
		firestoreClientMutex:    &sync.Mutex{},
	}
}

// Connect is the Firebase DB connection
func (p Production) getFirestoreConnection() (*firestore.Client, error) {
	// TODO mutex
	if p.firestoreClient != nil {
		return p.firestoreClient, nil
	}
	var client *firestore.Client
	ctx := context.Background()
	opt := option.WithCredentialsFile(p.config.FirebaseConfig)

	config := &firebase.Config{}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, err
	}
	client, err = app.Firestore(ctx)
	if err != nil {
		return nil, err

	}

	return client, nil
}
