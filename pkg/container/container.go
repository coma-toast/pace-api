package container

import "cloud.google.com/go/firestore"

type Container struct {
	Firebase *firestore.Client
}
