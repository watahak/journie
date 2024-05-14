package main

import (
	"context"
	"encoding/json"
	"fmt"
	chatsession "journie/pkg/chat-session"
	firebaseClient "journie/pkg/firebase"
	"journie/pkg/generative"
	"journie/pkg/messaging"
	"journie/pkg/pubsub"
	"journie/pkg/templates"
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"

	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	ctx := context.Background()

	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	enverr := godotenv.Load(curDir + "/.env")
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}

	// init google pubsub
	pubsuberr := pubsub.SubscribeToTopic(context.Background(), os.Getenv("FIREBASE_PROJECT_ID"), "remind-topic", "remind-sub")

	if pubsuberr != nil {
		fmt.Println("error with pubsub")
		log.Fatal(err)
	}

	// init gemini
	generative.Init()

	// init chat sessions
	chatsession.Init()

	// init telebot
	teleErr := messaging.Init()

	if teleErr != nil {
		log.Fatalln(teleErr)
		return
	}

	// init firebase and firestore
	firebaseErr := firebaseClient.Init(ctx)

	if firebaseErr != nil {
		log.Fatalln(firebaseErr)
	}

	messaging.TeleBot.Handle("/start", func(c tele.Context) error {
		var (
			username = c.Sender().Username
		)

		message := templates.WelcomeMessageSharedApiKey(username)

		return c.Send(message, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	})

	// to remove before live
	messaging.TeleBot.Handle("/summarize", func(c tele.Context) error {
		var (
			userId = int(c.Sender().ID)
		)

		cs, err := chatsession.ChatSessionClient.GetOrCreateChatSession(messaging.GetPlatformUserId(fmt.Sprint(userId)))

		if err != nil {
			log.Print(err)
			return c.Send("Error retrieving chat session")
		}

		var platformUserId = messaging.GetPlatformUserId(fmt.Sprint(userId))

		analysis, err := chatsession.IngestChatSession(cs, platformUserId)

		if err != nil {
			log.Print(err)
			return c.Send("Error retrieving chat session")
		}

		out, err := json.Marshal(analysis)

		fmt.Println(string(out))

		if err != nil {
			panic(err)
		}

		return c.Send(string(out))
	})

	messaging.TeleBot.Handle("/clear", func(c tele.Context) error {
		var (
			userId = int(c.Sender().ID)
		)

		err := chatsession.ChatSessionClient.DeleteChatSession(messaging.GetPlatformUserId(fmt.Sprint(userId)))

		if err != nil {
			return c.Send("Error deleting chat session")
		}

		return c.Send("Chat session deleted")
	})

	messaging.TeleBot.Handle(tele.OnText, func(c tele.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.

		var (
			sender = c.Sender()
			text   = c.Text()
		)

		userId := int(sender.ID)
		log.Println("userId:", userId)

		ctx := context.Background()

		// Initialize chat session
		// @todo close chat session if user is inactive
		cs, err := chatsession.ChatSessionClient.GetOrCreateChatSession(messaging.GetPlatformUserId(fmt.Sprint(userId)))

		if err != nil {
			log.Fatal(err)
			return c.Send("Error creating chat session")
		}

		messaging.TeleBot.Notify(sender, tele.Typing)

		resp, err := cs.SendMessage(ctx, genai.Text(text))

		// @todo, if err occurs due to safety, reflect in message, and recover the history by creating a new session
		if err != nil {
			log.Println("Error sending message to Gemini:", err)
			logrus.Errorf("An error occurred: %v", err)
			return c.Send("Error processing your request")
		}

		out, err := generative.ResponseToString(resp)
		if err != nil {
			return c.Send("Error generating chat response")
		}

		return c.Send(out)
	})

	messaging.TeleBot.Handle(tele.OnPhoto, func(c tele.Context) error {
		return c.Send(string("Sorry! I am unable to process images as of now!"))
	})

	messaging.TeleBot.Handle(tele.OnVideo, func(c tele.Context) error {
		return c.Send(string("Sorry! I am unable to process videos as of now!"))
	})

	messaging.TeleBot.Handle(tele.OnVoice, func(c tele.Context) error {
		return c.Send(string("Sorry! I am unable to process voice messages as of now!"))
	})

	messaging.TeleBot.Start()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
