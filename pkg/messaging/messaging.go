package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	chatsession "journie/pkg/chat-session"
	"journie/pkg/generative"
	"journie/pkg/templates"
	"journie/pkg/users"
	"journie/pkg/utility"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	tele "gopkg.in/telebot.v3"
)

var TeleBot *tele.Bot

var (
	// Universal markup builders.
	// menu     = &tele.ReplyMarkup{ResizeKeyboard: true} // located on keyboard
	selector = &tele.ReplyMarkup{} // located with message bubble

	// Reply buttons.
	// btnHelp = menu.Text("Help")
	// btnSettings = menu.Text("Settings")

	// Inline buttons.
	// btnWhy = selector.Data("Why do I need this?", "gemini-reason")
)

type UserModel struct {
	Platform string `json:"platform"`
	UserId   string `json:"userId"`
}

const (
	Telegram string = "telegram"
)

func Init() error {
	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	TeleBot, err := tele.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return err
	}

	TeleBot.Handle("/gemini_key", func(c tele.Context) error {
		return c.Send(templates.GeminiKeyInstructions, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: selector})
	})

	// // On reply button pressed (message)
	// TeleBot.Handle(&btnHelp, func(c tele.Context) error {
	// 	return c.Send("Here is some help: ...!")
	// })

	// // On inline button pressed (callback)
	// TeleBot.Handle(&btnWhy, func(c tele.Context) error {
	// 	return c.Send(templates.GeminiKeyReason, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	// })

	// handle default /start command from telegram
	TeleBot.Handle("/start", func(c tele.Context) error {
		var username = c.Sender().Username

		message := templates.WelcomeMessageSharedApiKey(username)

		return c.Send(message, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	})

	// handle manual summarize
	TeleBot.Handle("/summarize", func(c tele.Context) error {
		var userId = int(c.Sender().ID)
		platformUserId, err := GetPlatformUserId(fmt.Sprint(userId))
		if err != nil {
			log.Printf("Error handling user id %d: %v", userId, err)
			return c.Send("Error handling user id")
		}

		cs, err := chatsession.ChatSessionClient.GetOrCreateChatSession(platformUserId)
		if err != nil {
			log.Printf("Error retrieving or creating chat session: %v", err)
			return c.Send("Error retrieving chat session")
		}

		analysis, err := chatsession.IngestChatSession(cs, platformUserId)
		if err != nil {
			log.Printf("Error ingesting chat session: %v", err)
			return c.Send("Error retrieving chat session")
		}

		out, err := json.Marshal(analysis)
		if err != nil {
			log.Printf("Error marshalling analysis: %v", err)
			return c.Send("Error processing analysis")
		}

		return c.Send(string(out))
	})

	// handle manual command to clear session
	TeleBot.Handle("/clear", func(c tele.Context) error {
		var userId = int(c.Sender().ID)
		platformUserId, err := GetPlatformUserId(fmt.Sprint(userId))
		if err != nil {
			log.Printf("Error handling user id %d: %v", userId, err)
			return c.Send("Error handling user id")
		}

		err = chatsession.ChatSessionClient.DeleteChatSession(platformUserId)
		if err != nil {
			log.Printf("Error deleting chat session: %v", err)
			return c.Send("Error deleting chat session")
		}

		return c.Send("Chat session deleted")
	})

	// All other text messages to be handled by Journie
	TeleBot.Handle(tele.OnText, func(c tele.Context) error {
		ctx := context.Background()
		var (
			sender = c.Sender()
			text   = c.Text()
		)
		userId := int(sender.ID)
		platformUserId, err := GetPlatformUserId(fmt.Sprint(userId))
		if err != nil {
			log.Printf("Error handling user id %d: %v", userId, err)
			return c.Send("Error handling user id")
		}

		// Initialize chat session
		cs, err := chatsession.ChatSessionClient.GetOrCreateChatSession(platformUserId)
		if err != nil {
			log.Printf("Error creating chat session for user %d: %v", userId, err)
			return c.Send("Error creating chat session")
		}

		TeleBot.Notify(sender, tele.Typing)

		resp, err := cs.SendMessage(ctx, genai.Text(text))
		if err != nil {
			log.Printf("Error sending message to chat session for user %d: %v", userId, err)
			// @todo, if err occurs due to safety, reflect in message, and recover the history by creating a new session
			return c.Send("Error processing your request")
		}

		out, err := generative.ResponseToString(resp)
		if err != nil {
			log.Printf("Error generating chat response for user %d: %v", userId, err)
			return c.Send("Error generating chat response")
		}

		return c.Send(out)
	})

	TeleBot.Handle(tele.OnPhoto, func(c tele.Context) error {
		return c.Send(string("Sorry! I am unable to process images as of now!"))
	})

	TeleBot.Handle(tele.OnVideo, func(c tele.Context) error {
		return c.Send(string("Sorry! I am unable to process videos as of now!"))
	})

	TeleBot.Handle(tele.OnVoice, func(c tele.Context) error {
		return c.Send(string("Sorry! I am unable to process voice messages as of now!"))
	})

	TeleBot.Start()

	return nil
}

func GetPlatformUserId(userId string) (string, error) {
	if userId == "" {
		return "", fmt.Errorf("invalid user ID: expected non empty string")
	}
	return fmt.Sprintf("%s-%s", Telegram, userId), nil
}

func ParsePlatformUserId(platformUserId string) (*UserModel, error) {
	// Split the platform ID and user ID
	splitted := strings.SplitN(platformUserId, "-", 2) // Ensure at most one separator

	// Check for expected number of parts
	if len(splitted) != 2 {
		return nil, fmt.Errorf("invalid platform user ID format: expected '{platform}-userId'")
	}

	// Return parsed data as a map with string values
	return &UserModel{
		Platform: splitted[0],
		UserId:   splitted[1],
	}, nil
}

// RemindDaily triggered to send reminder messsage to users
// @todo handle timezone
func RemindDaily() {
	targetDate := time.Now().Format("2006-01-02")
	inactiveUsers := users.GetUsersWithoutSession(targetDate)

	// throttle, telegram rate limits ~30 per second
	throttleDuration := 100 * time.Millisecond
	throttle := utility.NewThrottle(throttleDuration)

	for _, userId := range inactiveUsers {
		teleUserId, err := ParsePlatformUserId(userId)
		if err != nil {
			log.Printf("Error parsing user ID %s: %v", userId, err)
			continue
		}

		var userIntValue int64
		_, converr := fmt.Sscan(teleUserId.UserId, &userIntValue) // Handle potential error with _ (blank identifier)
		if converr != nil {
			log.Printf("Error converting teleUserId.UserId %s to int: %v", userId, converr)
			continue
		}

		u := &tele.User{ID: userIntValue}

		go func() {
			throttle.Process()
			_, err := TeleBot.Send(u, "Hi, take 5 minutes to write a journal entry!")
			if err != nil {
				log.Printf("Error sending reminder to user %d: %v", u.ID, err)
			}
		}()
	}
}

// SummarizeDaily:
// 1. get users with last chat session for the day
// 2. summarize chat sessions (handle gemini rate limits @ ~60 per minute)
// 3. delete chat session once summarized
func SummarizeDaily() {
	// assuming day "ends" at 4am
	// floor hour and minutes, get date and minus 1 day
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	users := users.GetUsersWithSession(time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location()))

	// generate content. debounce and rate limit
	throttleDuration := 1000 * time.Millisecond
	throttle := utility.NewThrottle(throttleDuration)

	for _, userId := range users {
		chatSession := chatsession.ChatSessionClient.GetChatSession(userId)
		if chatSession == nil {
			log.Printf("Chat session not found for user %s", userId)
			continue
		}

		go func(platformUserId string) {
			throttle.Process()
			// Call the original function with converted arguments
			if _, err := chatsession.IngestChatSession(chatSession, platformUserId); err != nil {
				log.Printf("Error summarizing chat session for user %s: %v", platformUserId, err)
				return
			}
			if err := chatsession.ChatSessionClient.DeleteChatSession(platformUserId); err != nil {
				log.Printf("Error deleting chat session for user %s: %v", platformUserId, err)
			}
		}(userId)
	}
}

func testLog(userId string) error {
	time.Sleep(100 * time.Millisecond)
	// Call the original function with converted arguments
	log.Println("Handling:", userId)
	return nil
}

func TestLoop() {
	targetDate := time.Now().Format("2006-01-02")
	inactiveUsers := users.GetUsersWithoutSession(targetDate)

	throttleDuration := 5000 * time.Millisecond
	throttle := utility.NewThrottle(throttleDuration)

	for _, userId := range inactiveUsers {
		go func(id string) {
			throttle.Process()
			testLog(id)
		}(userId)
	}
}
