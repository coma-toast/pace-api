package container

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/coma-toast/pace-api/pkg/paceconfig"
	"github.com/coma-toast/pace-api/pkg/provider/contact"
	"github.com/coma-toast/pace-api/pkg/provider/user"
	"google.golang.org/api/option"
)

// Container exposes data providers
type Container interface {
	UserProvider() (user.Provider, error)
	ContactProvider() (contact.Provider, error)
}

// Production is our production container for our external connections
type Production struct {
	config *paceconfig.Config
	// Providers
	userProvider    *user.DatabaseProvider
	contactProvider *contact.DatabaseProvider
	// Clients
	firestoreClient *firestore.Client
	// Mutex Locks
	userProviderMutex    *sync.Mutex
	contactProviderMutex *sync.Mutex
	firestoreClientMutex *sync.Mutex
}

// UserProvider provides the user provider
func (p Production) UserProvider() (user.Provider, error) {
	if p.userProvider != nil {
		return p.userProvider, nil
	}
	// TODO: copy Hub user provider (mutex lock, etc)
	firestoreConnection, err := p.getFirestoreConnection()
	if err != nil {
		return nil, err
	}
	p.userProvider = &user.DatabaseProvider{
		Database: firestoreConnection,
	}

	return p.userProvider, nil
}

// ContactProvider provides the contact provider
func (p Production) ContactProvider() (contact.Provider, error) {
	if p.contactProvider != nil {
		return p.contactProvider, nil
	}
	// TODO: copy Hub contact provider (mutex lock, etc)
	firestoreConnection, err := p.getFirestoreConnection()
	if err != nil {
		return nil, err
	}
	p.contactProvider = &contact.DatabaseProvider{
		Database: firestoreConnection,
	}

	return p.contactProvider, nil
}

// NewProduction builds a container with all of the config
func NewProduction(paceconfig *paceconfig.Config) Container {
	return &Production{
		config:               paceconfig,
		userProviderMutex:    &sync.Mutex{},
		firestoreClientMutex: &sync.Mutex{},
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
