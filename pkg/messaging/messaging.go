package messaging

import (
	"fmt"
	chatsession "journie/pkg/chat-session"
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

	teleBot, err := tele.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return err
	}

	TeleBot = teleBot

	return nil
}

func GetPlatformUserId(userId string) string {
	return fmt.Sprintf("%s-%s", Telegram, userId)
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

// todo handle timezone
func RemindDaily() {
	targetDate := time.Now().Format("2006-01-02")
	inactiveUsers := users.GetUsersWithoutSession(targetDate)

	// debounce, telegram rate limits ~30 per second
	debounceDuration := 100 * time.Millisecond
	debouncedFn := utility.Debounce(debouncedSendMsg, debounceDuration)

	for _, userId := range inactiveUsers {

		teleUserId, err := ParsePlatformUserId(userId)

		if err != nil {
			log.Println(err)
			continue
		}

		var userIntValue int64

		_, converr := fmt.Sscan(teleUserId.UserId, &userIntValue) // Handle potential error with _ (blank identifier)

		if converr != nil {
			log.Println("Error converting teleUserId.UserId to int:", err)
			continue
		}

		u := &tele.User{ID: userIntValue}

		go debouncedFn(u, "Hi, take 5 minutes to write a journal entry!")
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
	debounceDuration := 1000 * time.Millisecond
	debouncedFn := utility.Debounce(debouncedIngest, debounceDuration)

	for _, userId := range users {

		chatSession := chatsession.ChatSessionClient.GetChatSession(userId)

		if chatSession == nil {
			continue
		}

		// Calling debouncedFn multiple times within debounceDuration will trigger only once after the duration elapses
		go debouncedFn(chatSession, userId)

	}
}

func debouncedIngest(args ...interface{}) {
	if len(args) != 2 {
		log.Println("debouncedIngest argument mismatch")
		return
	}
	chatSession, ok1 := args[0].(*genai.ChatSession)
	platformUserId, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		log.Println("debouncedIngest invalid argument types")
		return
	}
	// Call the original function with converted arguments
	chatsession.IngestChatSession(chatSession, platformUserId)

	// dirty delete call
	chatsession.ChatSessionClient.DeleteChatSession(platformUserId)
}

func debouncedSendMsg(args ...interface{}) {
	if len(args) != 2 {
		log.Println("debouncedSendMsg argument mismatch")
		return
	}
	to, ok1 := args[0].(tele.Recipient)
	what, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		log.Println("debouncedIngest invalid argument types")
		return
	}
	// Call the original function with converted arguments
	TeleBot.Send(to, what)
}
