package chatsession

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	firebaseClient "journie/pkg/firebase"
	"journie/pkg/generative"
	"log"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/generative-ai-go/genai"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"google.golang.org/api/iterator"
)

var ChatSessionClient *ChatSession

type ChatSession struct {
	Sessions map[string]*genai.ChatSession // Map of user IDs to Gemini clients
	mu       sync.Mutex                    // Mutex to synchronize access to the map
}

type AnalysisResult struct {
	Summary   string    `json:"summary"`
	Mood      []string  `json:"mood"`
	CreatedAt time.Time `json:"createdAt"`
}

func Init() {
	ChatSessionClient = &ChatSession{
		Sessions: make(map[string]*genai.ChatSession),
	}
}

func MapToAnalysisResult(data map[string]interface{}) (*AnalysisResult, error) {
	var result AnalysisResult
	err := mapstructure.Decode(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func AnalysisResultToHistory(data *AnalysisResult) string {
	return fmt.Sprintf("On the date %s, the following conversation happened with you and the user, where the user is in second-person: '%s'. User's mood was: %s",
		data.CreatedAt.Format("2006-01-02"), data.Summary, strings.Join(data.Mood, ", "))
}

// GetChatSession retrieves chat-session of user if it exists.
// userId should be a string in format {platform}-{indentifier}
func (cs *ChatSession) GetChatSession(userId string) *genai.ChatSession {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if session, ok := cs.Sessions[userId]; ok {
		return session
	}

	return nil
}

// GetOrCreateChatSession retrieves existing chat session or creates new one with historical context
func (cs *ChatSession) GetOrCreateChatSession(userID string) (*genai.ChatSession, error) {
	ctx := context.Background()

	existingSession := cs.GetChatSession(userID)
	if existingSession != nil {
		return existingSession, nil
	}

	model := generative.GetUserModel()
	chatSession := model.StartChat()
	fmt.Printf("New chat session created for user %s", userID)

	err := firebaseClient.UpsertUserLastCreatedSession(ctx, userID)
	if err != nil {
		log.Printf("Error updating last created session for user %s: %v", userID, err)
		return nil, err
	}

	// get past x days of summaries for user from firestore, insert into history
	userRef := firebaseClient.FirestoreClient.Collection("users").Doc(userID)
	subcollectionRef := userRef.Collection("entries")
	iter := subcollectionRef.OrderBy("createdAt", firestore.Asc).Limit(30).Documents(ctx)

	var summaries []AnalysisResult
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if doc == nil {
			break
		}

		result, err := MapToAnalysisResult(doc.Data())
		if err != nil {
			log.Printf("Error mapping document to AnalysisResult for user %s: %v", userID, err)
			return nil, err
		}

		summaries = append(summaries, *result)
	}

	if len(summaries) != 0 {
		histories := lo.Map(summaries, func(i AnalysisResult, index int) string {
			return AnalysisResultToHistory(&i)
		})

		parts := make([]genai.Part, len(histories))
		for i, history := range histories {
			parts[i] = genai.Text(history)
		}

		chatSession.History = []*genai.Content{
			{
				Parts: parts,
				Role:  "model",
			},
		}
	}

	cs.Sessions[userID] = chatSession
	return chatSession, nil
}

// DeleteChatSession delete chat session for user
func (cs *ChatSession) DeleteChatSession(userID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.Sessions, userID)

	return nil
}

// IngestChatSession summarize chat seesion for user
// Persists entry into db
func IngestChatSession(chatSession *genai.ChatSession, platformUserId string) (*AnalysisResult, error) {
	ctx := context.Background()

	response, err := generative.SummarizeSession(chatSession)
	if err != nil {
		log.Println("Error generating summary:", err)
		return nil, err
	}

	jsonString := response.Candidates[0].Content.Parts[0]
	text, ok := jsonString.(genai.Text)
	if !ok {
		return nil, errors.New("failed to convert to text")
	}

	cleanJSON := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(string(text), "``` \n"), "```json\n"))

	var result AnalysisResult
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	now := time.Now()
	result.CreatedAt = now

	// store text into firestore
	collectionPath := fmt.Sprintf("users/%s/%s", platformUserId, "entries")
	_, err = firebaseClient.FirestoreClient.Collection(collectionPath).Doc(now.Format("2006-01-02")).Set(ctx, map[string]interface{}{
		"summary":   &result.Summary,
		"mood":      &result.Mood,
		"createdAt": &result.CreatedAt,
	})

	if err != nil {
		return nil, fmt.Errorf("error saving summary to firestore: %w", err)
	}

	return &result, nil
}
