package llm

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	Client *openai.Client
}

func NewOpenAIClient() *OpenAIClient {
	apiKey := os.Getenv("OPENAI_API_KEY")

	// If not in environment, try loading from ~/.langmate.env
	if apiKey == "" {
		apiKey = loadAPIKeyFromConfig()
	}

	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not found. Set it in environment or ~/.langmate.env")
		fmt.Println("Example: echo 'OPENAI_API_KEY=sk-...' > ~/.langmate.env")
		os.Exit(1)
	}

	client := openai.NewClient(apiKey)
	return &OpenAIClient{Client: client}
}

func loadAPIKeyFromConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	configPath := filepath.Join(home, ".langmate.env")
	file, err := os.Open(configPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "OPENAI_API_KEY=") {
			return strings.TrimPrefix(line, "OPENAI_API_KEY=")
		}
	}
	return ""
}

func (c *OpenAIClient) TransferText(text string, lang string) (string, error) {
	prompt := fmt.Sprintf("You will receive a text and a destination language. "+
		"If the text is not in the destination language, translate it."+
		"If the text is already in the destination language, rephrase it for clarity and style."+
		"If the text is a single word or slang, please explain it in a way dictionary does,"+
		"When the input is word or slang, in the bottom please list some alternatives if there are any"+
		"\n\nInput Text: \" %s \"\nDestination Language: \" %s \"\n\nOutput:", text, lang)

	resp, err := c.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %w", err)
	}
	rst := resp.Choices[0].Message.Content
	return rst, nil
}

func (c *OpenAIClient) RephraseText(text string, lang string) (string, error) {
	resp, err := c.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a writing assistant. Rephrase the given text for clarity and better style. Output ONLY the rephrased text, nothing else. No explanations, no alternatives, no quotes around the text. Preserve the original meaning and tone.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %w", err)
	}
	return resp.Choices[0].Message.Content, nil
}
