package container

import "cloud.google.com/go/firestore"

// Make interface that contains other interfaces
type Container struct {
	Firestore *firestore.Client
}
