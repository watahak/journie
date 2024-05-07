package pubsub

import (
	"context"
	"fmt"
	"journie/pkg/messaging"
	"time"

	"cloud.google.com/go/pubsub"
)

// SubscribeToTopic subscribes to a Pub/Sub topic and receives messages.
func SubscribeToTopic(ctx context.Context, projectID, topicName string, subscriptionName string) error {
	client, err := pubsub.NewClient(ctx, projectID)
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
			fmt.Println("Received message:", string(msg.Data))
			// Process the message data and acknowledge it

			// hourly handler

			var now = time.Now()

			// if now.UTC().Hour() == 14 { // sg 10pm
			go messaging.RemindDaily()

			if now.UTC().Hour() == 16 { // sg 12am
				messaging.SummarizeDaily()
			}

			msg.Ack()
		})
		if err != nil {
			fmt.Println("Error receiving messages:", err)
		}
	}()

	return nil
}
