package pubsub

import (
	"context"
	"fmt"
	"journie/pkg/messaging"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// SubscribeToTopic subscribes to a Pub/Sub topic and receives messages.
func SubscribeToTopic(ctx context.Context, projectID, topicName string, subscriptionName string) error {
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS")))

	client, err := pubsub.NewClient(ctx, projectID, opt)
	if err != nil {
		fmt.Println("NewClient failed")
		return err
	}

	// Reference the topic
	// topic := client.Topic(topicName)

	// Create a subscription (or use an existing one)
	sub := client.Subscription(subscriptionName)

	go func() {
		// Receive messages concurrently
		err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			fmt.Println("Received message at:", msg.PublishTime)

			// hourly handler
			var now = msg.PublishTime

			// go messaging.TestLoop()

			if now.UTC().Hour() == 14 { // sg 10pm
				go messaging.RemindDaily()
			}

			if now.UTC().Hour() == 20 { // sg 4am
				go messaging.SummarizeDaily()
			}

			msg.Ack()
		})
		if err != nil {
			fmt.Println("Error receiving messages:", err)
		}
	}()

	return nil
}
