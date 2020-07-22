package user

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/coma-toast/pace-api/pkg/entity"
)

// DatabaseProvider is a user.Provider the uses a database
type DatabaseProvider struct {
	Database *firestore.Client
}

// ErrUserNotFound if no users are found
var ErrUserNotFound = errors.New("User not found")

// GetByUsername gets a User by username
func (d *DatabaseProvider) GetByUsername(username string) (entity.User, error) {
	return d.getByUsername(username)
}

func (d *DatabaseProvider) getByUsername(username string) (entity.User, error) {
	var user entity.User

	users := d.Database.Collection("users").Where("username", "==", username).Documents(context.TODO())
	allMatchingUsers, err := users.GetAll()
	if err != nil {
		return entity.User{}, err
	}

	for _, fbUser := range allMatchingUsers {
		err = fbUser.DataTo(&user)
		if err != nil {
			return entity.User{}, fmt.Errorf("ERROR: User error - Firestore.DataTo() error %w, for user %s", err, username)
		}
		return user, nil
		// data = append(data, fbUser.Data())
	}

	return entity.User{}, ErrUserNotFound
}
