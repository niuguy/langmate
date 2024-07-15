package llm

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	Client *openai.Client
}

func NewOpenAIClient() *OpenAIClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set.")
		os.Exit(1)
	}

	client := openai.NewClient(apiKey)
	return &OpenAIClient{Client: client}
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
