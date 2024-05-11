package generative

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var GenAiClient *GenAiManager

type GenAiManager struct {
	Client *genai.Client
	Model  *genai.GenerativeModel
}

func Init() error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel(os.Getenv("GEMINI_MODEL"))

	if err != nil {
		return err
	}

	GenAiClient = &GenAiManager{
		Client: client,
		Model:  model,
	}

	fmt.Println("gemini client created")

	return nil
}

func ResponseToString(resp *genai.GenerateContentResponse) (string, error) {
	part := resp.Candidates[0].Content.Parts[0]

	text, ok := part.(genai.Text)
	if !ok {
		return "", errors.New("unexpected part type")
	}
	return string(text), nil
}
