package messaging

import (
	"fmt"
	chatsession "journie/pkg/chat-session"
	"journie/pkg/templates"
	"journie/pkg/users"
	"journie/pkg/utility"
	"log"
	"os"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

var TeleBot *tele.Bot

var (
	// Universal markup builders.
	menu     = &tele.ReplyMarkup{ResizeKeyboard: true} // located on keyboard
	selector = &tele.ReplyMarkup{}                     // located with message bubble

	// Reply buttons.
	btnHelp     = menu.Text("Help")
	btnSettings = menu.Text("Settings")

	// Inline buttons.
	//
	// Pressing it will cause the client to
	// send the bot a callback.
	//
	// Make sure Unique stays unique as per button kind
	// since it's required for callback routing to work.
	//
	btnWhy = selector.Data("Why do I need this?", "gemini-reason")
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

	teleBot, err := tele.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return err
	}

	TeleBot = teleBot

	// inline keybpard
	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)
	selector.Inline(
		selector.Row(btnWhy),
	)

	TeleBot.Handle("/gemini_key", func(c tele.Context) error {
		return c.Send(templates.GeminiKeyInstructions, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: selector})
	})

	// On reply button pressed (message)
	TeleBot.Handle(&btnHelp, func(c tele.Context) error {
		return c.Send("Here is some help: ...!")
		// return c.Edit("Here is some help: ...")
	})

	// On inline button pressed (callback)
	TeleBot.Handle(&btnWhy, func(c tele.Context) error {
		return c.Send(templates.GeminiKeyReason, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	})

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

	// throttle, telegram rate limits ~30 per second
	throttleDuration := 100 * time.Millisecond
	throttle := utility.NewThrottle(throttleDuration)

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

		go func() {
			throttle.Process()
			TeleBot.Send(u, "Hi, take 5 minutes to write a journal entry!")
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
			continue
		}

		go func(platformUserId string) {
			throttle.Process()
			// Call the original function with converted arguments
			chatsession.IngestChatSession(chatSession, platformUserId)
			chatsession.ChatSessionClient.DeleteChatSession(platformUserId)
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
