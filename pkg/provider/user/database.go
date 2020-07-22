package user

import (
	"cloud.google.com/go/firestore"
	"github.com/coma-toast/pace-api/pkg/entity"
)

// DatabaseProvider is a user.Provider the uses a database
type DatabaseProvider struct {
	Database *firestore.Client
}

// GetByUsername gets a User by username
func (d *DatabaseProvider) GetByUsername(username string) (entity.User, error) {
	return d.getByUsername(username)
}

func (d *DatabaseProvider) getByUsername(username string) (entity.User, error) {
	var user entity.User

	d.Database.Connect().GetAll()
	// users := d.Database.Collection("users").Where("username", "==", username).Documents(ctx)
	// allMatchingUsers, err := users.GetAll()

	// if err != nil {
	// 	rollbar.Warning(
	// 		fmt.Sprintf("Error getting user %s from Firebase: %e", username, err))
	// }
	// for _, fbUser := range allMatchingUsers {
	// 	data = append(data, fbUser.Data())
	// }

	// encoder := json.NewEncoder(w)
	// if err := encoder.Encode(&data); err != nil {
	// 	rollbar.Warning(fmt.Sprintf("Error encoding JSON when retrieving a User: %e", err))
	// }
	// log.Println(data)
	// //

	// if err := d.Database.Where("`uuid` = ?", uuid).First(&dbUser).Error; err != nil {
	// 	return entity.User{}, err
	// }

	return user, nil
}
