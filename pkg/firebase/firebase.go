package firebaseClient

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client

func Init(ctx context.Context) error {
	conf := &firebase.Config{ProjectID: os.Getenv("FIREBASE_PROJECT_ID")}
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIAL_PATH"))

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	FirestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func UpdateUserLastCreatedSession(ctx context.Context, platformUserId string) error {

	var now = time.Now()

	_, err := FirestoreClient.Collection("users").Doc(platformUserId).Update(ctx, []firestore.Update{
		{
			Path:  "lastCreatedSession",
			Value: now,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
