package users

import (
	"context"
	"fmt"
	firebaseClient "journie/pkg/firebase"
	"time"

	"google.golang.org/api/iterator"
)

func GetUsersWithoutSession(datestring string) []string {
	ctx := context.Background()

	// Create a reference to the users collection
	client := firebaseClient.FirestoreClient

	now := time.Now()

	// Create a new time object with the desired date and zeroed hours and minutes
	currentDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	iter := client.Collection("users").Where(
		"lastCreatedSession", "<", currentDay,
	).Documents(ctx)

	var usersWithoutSession []string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			fmt.Println("Retrieved all matching users")
			break
		}
		if err != nil {
			fmt.Println("Error getting document:", err)
			continue
		}

		// Process the user document data (doc.Data())
		usersWithoutSession = append(usersWithoutSession, doc.Ref.ID)
	}

	userCount := len(usersWithoutSession)
	fmt.Printf("Found %d users without session\n", userCount)

	return usersWithoutSession
}
