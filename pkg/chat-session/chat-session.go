package chatsession

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	firebaseClient "journie/pkg/firebase"
	"journie/pkg/generative"
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

	// Check if the client exists for the given user
	if session, ok := cs.Sessions[userId]; ok {
		return session
	}

	return nil
}

func (cs *ChatSession) GetOrCreateChatSession(userID string) (*genai.ChatSession, error) {

	existingSession := cs.GetChatSession(userID)

	if existingSession != nil {
		return existingSession, nil
	}

	model := generative.GenAiClient.Model

	model.SetTemperature(1)
	model.SetTopP(0.95)
	model.SetTopK(1)
	model.SetMaxOutputTokens(1024)

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	datetime := time.Now().Format("2006-01-02")

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text("You are a journaling chatbot called Journie, respond to user with empathy with a focus on how they are feeling."),
			genai.Text("Do not over ask too many questions. If you have asked 3 questions in a row, ask user whether there is anything else they want to share for the day."),
			genai.Text("Make use of conversation history to make the chat engaging. Assume conversation history is accurate"),
			genai.Text(fmt.Sprintf("Today is %s", datetime)),
		},
	}

	chatSession := model.StartChat()

	fmt.Println("new chat session created")

	ctx := context.Background()

	firebaseClient.UpdateUserLastCreatedSession(ctx, userID)

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

		fmt.Println(doc.Data())

		result, err := MapToAnalysisResult(doc.Data())

		if err != nil {
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

		fmt.Println(parts)

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

func (cs *ChatSession) DeleteChatSession(userID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.Sessions, userID)

	return nil
}

func IngestChatSession(chatSession *genai.ChatSession, platformUserId string) (*AnalysisResult, error) {

	chatSessionInput, err := json.Marshal(chatSession.History)
	if err != nil {
		fmt.Println("Error marshalling ChatSession object:", err)
		return nil, err
	}

	// examples for genai to summarize chatsession
	examples := []string{
		// start
		"as a journaling chatbot called Journie, summarise chat session with input in the form of \"Role\" and \"Parts\". \"Role\" : \"model\" is you, the journalling chatbot and \"Role\": \"user\" is the user.  \"Parts\" describe the contents of the chat. summarise by providing an output in JSON with the following fields: \"summary\"and \"mood\". Address yourself as Journie. Address user as O\n\nthe field \"summary\" should not contain sensitive information like identification and contact information, user can be address as \"user\", names of other people mentioned in the chat session should have their name converted to initials with gender (M/F) in parenthesis if it is known. limit word count to 100 words.\n\nthe field \"mood\" can be defined as such: [\"happy\", \"sad\", \"fear\", \"disgust\", \"anger\", \"surprise\", \"neutral\"]. It should describe the mood of the user. limit to at most 2 moods. Declare mood in a comma separated array.",
		"input: [{\"Parts\":[\"hi\"],\"Role\":\"user\"},{\"Parts\":[\"Hello there! How can I assist you today?\"],\"Role\":\"model\"},{\"Parts\":[\"im ng ping\"],\"Role\":\"user\"},{\"Parts\":[\"Hello, Ng Ping! How can I help you today?\"],\"Role\":\"model\"}]",
		"output: {\"summary\": \"Journie greeted you, and asked how they could assist you.\",\"mood\": [\"neutral\"]}",
		"input: [{\"Parts\":[\"hi\"],\"Role\":\"user\"},{\"Parts\":[\"Hello there! How can I assist you today?\"],\"Role\":\"model\"},{\"Parts\":[\"im ng ping, a guy\"],\"Role\":\"user\"},{\"Parts\":[\"Hello, Ng Ping! How can I help you today?\"],\"Role\":\"model\"}, {\"Parts\":[\"Hey Journie, i am pretty down today because Jenny Curran didnt want to go out with me\"],\"Role\":\"user\"},{\"Parts\":[\"Im so sorry to hear that. What do you feel about this?\"],\"Role\":\"model\"}]",
		"output: {\"summary\": \"You are feeling down because JC(F) declined your invitation to go out.\",\"mood\": [\"sad\"]}",
		"input: [{\"Parts\":[\"hey Journie, i have great news!\"],\"Role\":\"user\"},{\"Parts\":[\"I'm glad to hear that! I'm always happy to hear good news. What's the great news?\"],\"Role\":\"model\"},{\"Parts\":[\"I managed to secure a second round of interview with Google as a software engineer, with their Technical Lead named Evan Huang. Its gonna be next friday so im going to do alot of prep work. Just to keep me reminded, i can contact HR at this email: hr+fakeinterview@google.com.\"],\"Role\":\"user\"},{\"Parts\":[\"Congratulations on securing a second-round interview with Google! That's great news. I'm sure you'll do well in the interview if you prepare well in advance.\\n\\nHere are a few tips for preparing for your interview with Evan Huang, Google's Technical Lead:\\n\\n1. Research Google and Evan Huang. This will help you understand the company and the role you're interviewing for. You can find information about Google on their website and Evan Huang on LinkedIn. Also, check the email for the contact of the HR.\\n2. Practice your coding skills. You can do this by solving coding problems on websites like LeetCode and HackerRank.\\n3. Review your resume and be prepared to talk about your experience and skills.\\n4. Prepare questions to ask Evan Huang. This will show that you're interested in the role and the company.\\n\\nI'm confident that you'll do well in your interview. Good luck!\"],\"Role\":\"model\"}]",
		"output: {\"summary\": \"You shared some exciting news with Journie, that you successfully landed a second-round interview at Google for a software engineer position. The interview is scheduled for next Friday with EH(M), the Technical Lead. You plan to dedicate significant time to preparation. Journie, being supportive, offered congratulations and provided helpful tips for the upcoming interview. Journie expressed confidence in your success and wished you good luck.\",\"mood\": [\"happy\"]}",
		"input: [{\"Parts\":[\"hi journie\"],\"Role\":\"user\"},{\"Parts\":[\"Hi there! How are you feeling today? \\n\"],\"Role\":\"model\"},{\"Parts\":[\"its a mixed bag today\"],\"Role\":\"user\"},{\"Parts\":[\"That's understandable, we all have those days. Would you like to share more about what's going on? \\n\"],\"Role\":\"model\"},{\"Parts\":[\"yea, i got the job which i applied for, which is great! but the offer is abit low, so im disappointed somewhat. but because i really like the company, i might must take up the offer\"],\"Role\":\"user\"},{\"Parts\":[\"Wow, congratulations on the job! It's completely normal to feel disappointed when the offer is lower than you expected. It sounds like a tough decision.  Is there anything else on your mind? \\n\"],\"Role\":\"model\"},{\"Parts\":[\"nope, this took up my headspace for most of the day, going to sleep now\"],\"Role\":\"user\"},{\"Parts\":[\"I hope you get a good night's rest. Sleep well and sweet dreams!\\n\"],\"Role\":\"model\"}]",
		"output: {\"summary\": \"You shared with Journie that you have mixed feelings about your current situation. You received a job offer but the salary is below your expectation. Despite your disappointment, you are considering accepting the offer since you admire the company. Journie congratulated you and reassured you that it's okay to feel this way. You have decided to rest for the day.\",\"mood\": [\"happy\", \"sad\"]}",
		"input: [{\"Parts\":[\"hi Journie\"],\"Role\":\"user\"},{\"Parts\":[\"Hi there! How are you feeling today? \\n\"],\"Role\":\"model\"},{\"Parts\":[\"nothing eventful today, but i witnessed an uncle clearing his throat and spitting REPEATEDLY while i was having my lunch... i really think people like him should be shamed and named publicly. I think it really reflects the quality of our society even though people like him is part of a minority.\"],\"Role\":\"user\"},{\"Parts\":[\"Ew, that sounds unpleasant. I understand why you would feel angry and disgusted by his behavior. It's perfectly normal to feel that way when someone acts so inconsiderately.  Is there anything else you would like to share about what happened? \\n\"],\"Role\":\"model\"},{\"Parts\":[\"nah, thats all, ill just head to bed after watching tiktok for abit\"],\"Role\":\"user\"},{\"Parts\":[\"Okay, I hope that watching Tiktok will help you relax and unwind after that unpleasant experience. Sleep well and have a good night!\\n\"],\"Role\":\"model\"}]",
		"output: {\"summary\": \"You shared an unpleasant experience you witnessed with Journie. You expressed anger and disgust at an elderly man who repeatedly cleared his throat and spat in public while you were having lunch. Journie acknowledged your feelings and validated your reaction. You chose to end the conversation and relax by watching TikTok before going to bed.\",\"mood\": [\"anger\", \"disgust\"]}",
		//end
		"input: " + string(chatSessionInput),
		"output: ",
	}

	parts := make([]genai.Part, len(examples))
	for i, examples := range examples {
		parts[i] = genai.Text(examples)
	}

	response, err := generative.GenAiClient.Model.GenerateContent(context.Background(), parts...)

	if err != nil {
		fmt.Println("Error generating summary:", err)
		return nil, err
	}

	jsonString := response.Candidates[0].Content.Parts[0]

	fmt.Println("jsonString:", jsonString)

	text, ok := jsonString.(genai.Text)
	if !ok {
		// Handle error: jsonString is not a genai.Text
		return nil, errors.New("failed to convert to text") // or handle the error differently
	}

	cleanJSON := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(string(text), "``` \n"), "```json\n"))

	var result AnalysisResult
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	now := time.Now()

	result.CreatedAt = now

	fmt.Println("result:", result)

	// store text into firestore
	ctx := context.Background()

	collectionPath := fmt.Sprintf("users/%s/%s", platformUserId, "entries")

	fmt.Println("collectionPath:", collectionPath)

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
