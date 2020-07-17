package firebase

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

// Connect is the Firebase DB connection
func Connect(firebaseConfig string) *firestore.Client {
	var client *firestore.Client
	ctx := context.Background()
	opt := option.WithCredentialsFile(firebaseConfig)
	config := &firebase.Config{ProjectID: "pace-37aef"}
	spew.Dump(opt)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
	return client
}

// app, err := firebase.NewApp(context.Background(), nil, opt)
// if err != nil {
//   return nil, fmt.Errorf("error initializing app: %v", err)
// }

// https://medium.com/google-cloud/firebase-developing-serverless-functions-in-go-963cb011265d
// import (
// 	"context"
// 	"log"

// 	firebase "firebase.google.com/go"
// 	"firebase.google.com/go/db"
// )

// var client *db.Client

// func init() {
// 	ctx := context.Background()
// 	conf := &firebase.Config{
// 		DatabaseURL: "https://<CHANGE_ME>.firebaseio.com/",
// 	}
// 	app, err := firebase.NewApp(ctx, conf)
// 	if err != nil {
// 		log.Fatalf("firebase.NewApp: %v", err)
// 	}
// 	client, err = app.Database(ctx)
// 	if err != nil {
// 		log.Fatalf("app.Firestore: %v", err)
// 	}
// }
