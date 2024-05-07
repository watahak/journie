package messaging

import (
	"fmt"
	"journie/pkg/users"
	"log"
	"os"
	"strings"
	"time"

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

		TeleBot.Send(u, "Hi, take 5 minutes to write a journal entry!")
	}
}

func SummarizeDaily() {
	// get chat session
	// summarize chat sessions
	// delete chat sessions
}
