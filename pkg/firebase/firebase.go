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
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS")))

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
		return err
	}

	FirestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v", err)
		return err
	}

	return nil
}

func UpsertUserLastCreatedSession(ctx context.Context, platformUserId string) error {
	now := time.Now()
	_, err := FirestoreClient.Collection("users").Doc(platformUserId).Set(ctx, map[string]interface{}{
		"lastCreatedSession": now,
	}, firestore.MergeAll)

	return err
}
