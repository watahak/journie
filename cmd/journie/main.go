package main

import (
	"context"
	chatsession "journie/pkg/chat-session"
	firebaseClient "journie/pkg/firebase"
	"journie/pkg/generative"
	"journie/pkg/messaging"
	"journie/pkg/pubsub"
	"log"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go func() {
		if err := r.Run(":" + os.Getenv("PORT")); err != nil {
			log.Fatal(err)
		}
	}()

	ctx := context.Background()

	// init google pubsub
	pubsuberr := pubsub.SubscribeToTopic(context.Background(), os.Getenv("FIREBASE_PROJECT_ID"), "remind-topic", "remind-sub")
	if pubsuberr != nil {
		log.Fatal(pubsuberr)
	}

	// init gemini
	generative.Init()

	// init chat sessions
	chatsession.Init()

	// init firebase and firestore
	firebaseErr := firebaseClient.Init(ctx)
	if firebaseErr != nil {
		log.Fatal(firebaseErr)
	}

	// init telebot
	teleErr := messaging.Init()
	if teleErr != nil {
		log.Fatal(teleErr)
	}
}
