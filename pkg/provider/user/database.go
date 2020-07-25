package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/coma-toast/pace-api/pkg/entity"
	helper "github.com/coma-toast/pace-api/pkg/utils"
	"github.com/google/uuid"
	"github.com/rollbar/rollbar-go"
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

// GetAll gets a User by username
func (d *DatabaseProvider) GetAll() ([]entity.User, error) {
	return d.getAll()
}

func (d *DatabaseProvider) getAll() ([]entity.User, error) {
	var users []entity.User

	allUserData, err := d.Database.Collection("users").Documents(context.TODO()).GetAll()
	if err != nil {
		return []entity.User{}, err
	}

	for _, userData := range allUserData {
		var user entity.User
		err := userData.DataTo(&user)
		if err != nil {
			return []entity.User{}, fmt.Errorf("ERROR: GetAll(): Firestore.DataTo() error %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// AddUser is to update a user record
func (d *DatabaseProvider) AddUser(newUserData entity.User) (entity.User, error) {
	userRef, err := d.addUser(newUserData)
	if err != nil {
		return entity.User{}, err
	}
	rollbar.Info(fmt.Sprintf("Adding new user %s", newUserData.Username))
	updatedUserData, err := d.getByUserID(userRef.ID)
	if err != nil {
		return entity.User{}, err
	}

	return updatedUserData, nil
}

// UpdateUser is to update a user record
func (d *DatabaseProvider) UpdateUser(newUserData entity.UpdateUserRequest) (entity.User, error) {
	currentUserData, err := d.getByUsername(newUserData.Username)
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

	err = d.setByUserID(currentUserData.ID, updatedUser)
	if err != nil {
		return entity.User{}, err
	}
	updatedUserData, err := d.getByUserID(updatedUser.ID)
	if err != nil {
		return entity.User{}, err
	}

	return updatedUserData, nil
}

func (d *DatabaseProvider) addUser(userData entity.User) (entity.User, error) {
	newUUID := uuid.New().String()
	newUserData := entity.User{
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
	addUserResult, err := d.Database.Collection("users").Doc(newUUID).Set(context.TODO(), newUserData)
	if err != nil {
		return entity.User{}, fmt.Errorf("Error setting user %s by ID: %s", newUserData.Username, err)
	}
	rollbar.Info(fmt.Sprintf("User %s added at %s.", newUserData.Username, addUserResult))

	newUser, err := d.getByUserID(newUserData.ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("Error getting newly created user %s by ID: %s", newUserData.Username, err)
	}

	return newUser, nil
}

func (d *DatabaseProvider) getByUserID(userID string) (entity.User, error) {
	var user entity.User
	userData, err := d.Database.Collection("users").Doc(userID).Get(context.TODO())
	if err != nil {
		return entity.User{}, fmt.Errorf("Error getting user %s by ID: %s", userID, err)
	}
	userData.DataTo(&user)

	return user, nil
}

func (d *DatabaseProvider) setByUserID(userID string, userData entity.User) error {
	_, err := d.Database.Collection("users").Doc(userID).Set(context.TODO(), userData)
	if err != nil {
		return fmt.Errorf("Error setting user %s by ID: %s", userID, err)
	}

	return nil
}

func (d *DatabaseProvider) getByUsername(username string) (entity.User, error) {
	var user entity.User

	users := d.Database.Collection("users").Where("Username", "==", username).Documents(context.TODO())
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
