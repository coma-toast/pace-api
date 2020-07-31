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
	// * HERE
	// Aaron Ellington:pizza:  3:31 PM
	// You can have a map of entity.User
	// Then just add a method to that map type to convert it to a regular struct.	//
	var users []entity.User
	// userList := make(map[string]interface{})
	err := d.SharedProvider.GetAll(&users)
	// err := d.SharedProvider.GetAll(&userList)
	if err != nil {
		return nil, err
	}
	// mapstructure.Decode(userList, &users)

	return users, nil
}

// GetByUsername gets a User by username
func (d *DatabaseProvider) GetByUsername(username string) (entity.User, error) {
	var user entity.User

	err := d.SharedProvider.GetFirstBy("username", "==", username, &user)
	if err != nil {
		return entity.User{}, ErrUserNotFound
	}

	return user, nil
}

// AddUser is to update a user record
func (d *DatabaseProvider) AddUser(userData entity.User) (entity.User, error) {
	rollbar.Info(fmt.Sprintf("Adding new User to DB %s %s - %s", userData.FirstName, userData.LastName, userData.Username))

	existingUser, _ := d.GetByUsername(userData.Username)
	if (entity.User{}) != existingUser {
		return entity.User{}, fmt.Errorf("Error adding user %s: Username already exists", userData.Username)
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
	err := d.SharedProvider.Set(userData.ID, userData)
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

// UpdateUser is to update a user record
func (d *DatabaseProvider) UpdateUser(newUserData entity.UpdateUserRequest) (entity.User, error) {
	currentUserData, err := d.GetByUsername(newUserData.Username)
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

	return updatedUserData, nil
}

// DeleteUser is to update a user record
func (d *DatabaseProvider) DeleteUser(user entity.UpdateUserRequest) error {
	currentUserData, err := d.GetByUsername(user.Username)
	if err != nil {
		return err
	}
	err = d.SharedProvider.Delete(currentUserData.ID)
	if err != nil {
		return err
	}
	rollbar.Info(fmt.Sprintf("Deleted user %s", user.Username))

	return nil
}
