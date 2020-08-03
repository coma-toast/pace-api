package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/coma-toast/pace-api/pkg/entity"
	"github.com/coma-toast/pace-api/pkg/provider/firestoredb"
	helper "github.com/coma-toast/pace-api/pkg/utils"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
)

// DatabaseProvider is a user.Provider the uses a database
type DatabaseProvider struct {
	SharedProvider *firestoredb.DatabaseProvider
}

// ErrUserNotFound if no users are found
var ErrUserNotFound = errors.New("User not found")

// GetAll gets a User by username
func (d *DatabaseProvider) GetAll() ([]entity.User, error) {
	var users []entity.User
	err := d.SharedProvider.GetAll(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetByUsername gets a User by username
func (d *DatabaseProvider) GetByUsername(username string) (entity.User, error) {
	var user entity.User
	err := d.SharedProvider.GetFirstBy("Username", "==", username, &user)
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", err, ErrUserNotFound)
	}

	return user, nil
}

// Add is to update a user record
func (d *DatabaseProvider) Add(userData entity.User) (entity.User, error) {
	rollbar.Info(fmt.Sprintf("Adding new User to DB %s %s - %s", userData.FirstName, userData.LastName, userData.Username))

	var existingUser entity.User
	err := d.SharedProvider.GetFirstBy("Username", "==", userData.Username, &existingUser)
	if (entity.User{}) != existingUser {
		return entity.User{}, fmt.Errorf("Error adding user %s: Username already exists. ID: %s", userData.Username, existingUser.ID)
	}

	newUUID := uuid.New().String()
	userData = entity.User{
		ID:        newUUID,
		Created:   time.Now().String(),
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Role:      userData.Role,
		Username:  userData.Username,
		Password:  helper.Hash(userData.Password, newUUID),
		Email:     userData.Email,
		Phone:     userData.Phone,
		TimeZone:  userData.TimeZone,
		DarkMode:  userData.DarkMode,
	}
	err = d.SharedProvider.Set(userData.ID, userData)
	if err != nil {
		return entity.User{}, fmt.Errorf("Error setting user %s by ID: %s", userData.Username, err)
	}

	var newUser = entity.User{}
	err = d.SharedProvider.GetByID(userData.ID, &newUser)
	if err != nil {
		return entity.User{}, fmt.Errorf("Error getting newly created user %s by ID: %s", userData.Username, err)
	}

	rollbar.Info(fmt.Sprintf("User %s added.", userData.Username))
	return newUser, nil
}

// Update is to update a user record
func (d *DatabaseProvider) Update(newUserData entity.UpdateUserRequest) (entity.User, error) {
	var currentUserData entity.User
	err := d.SharedProvider.GetFirstBy("Username", "==", newUserData.Username, &currentUserData)
	if err != nil {
		return entity.User{}, err
	}

	rollbar.Info(fmt.Sprintf("Updating userID %s. \nOld Data: %v \nNew Data: %v", currentUserData.ID, currentUserData, newUserData))
	updatedUser := entity.User{
		ID:        currentUserData.ID,
		Created:   currentUserData.Created,
		FirstName: newUserData.FirstName,
		LastName:  newUserData.LastName,
		Role:      newUserData.Role,
		Username:  newUserData.Username,
		Password:  currentUserData.Password,
		Email:     newUserData.Email,
		Phone:     newUserData.Phone,
		TimeZone:  newUserData.TimeZone,
		DarkMode:  newUserData.DarkMode,
	}

	err = d.SharedProvider.Set(currentUserData.ID, updatedUser)
	if err != nil {
		return entity.User{}, err
	}

	var updatedUserData = entity.User{}
	err = d.SharedProvider.GetByID(currentUserData.ID, &updatedUserData)
	if err != nil {
		return entity.User{}, err
	}

	rollbar.Info(fmt.Sprintf("User %s updated.", updatedUserData.Username))

	return updatedUserData, nil
}

// Delete deletes a user
func (d *DatabaseProvider) Delete(user entity.User) error {
	rollbar.Info(fmt.Sprintf("Deleting User from DB: %s %s (%s)", user.FirstName, user.LastName, user.Username))

	var currentUser entity.User

	err := d.SharedProvider.GetByID(user.ID, &currentUser)
	if (entity.User{}) == currentUser {
		return fmt.Errorf("User not found")
	}

	err = d.SharedProvider.Delete(user.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted user %s", user.Username))

	return nil
}
