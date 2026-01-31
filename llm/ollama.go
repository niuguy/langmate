package llm

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

type OllamaClient struct {
	Client *api.Client
}

func NewOllamaClient() *OllamaClient {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}

	return &OllamaClient{Client: client}
}

func (c *OllamaClient) TransferText(text string, lang string) (string, error) {

	prompt := fmt.Sprintf("You will be given a text and a specified destination language. Follow these instructions based on the input:"+
		"\n1. If the text is not in the destination language, translate it."+
		"\n2. If the text is in the destination language, enhance its clarity and style."+
		"\n3. If the text consists of a single word or slang, define it as a dictionary would, and list any synonyms or related terms at the end."+
		"\n\nInput Text: \"%s\"\nDestination Language: \"%s\"\n\n", text, lang)

	req := &api.GenerateRequest{
		Model:  "llama3",
		Prompt: prompt,

		// set streaming to false
		Stream: new(bool),
	}

	ctx := context.Background()
	rspText := ""
	respFunc := func(resp api.GenerateResponse) error {
		rspText = resp.Response
		return nil
	}

	err := c.Client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}
	return rspText, nil

}

func (c *OllamaClient) RephraseText(text string, lang string) (string, error) {
	prompt := "You are a writing assistant. Rephrase the following text for clarity and better style. " +
		"Output ONLY the rephrased text, nothing else. No explanations, no alternatives, no quotes around the text. " +
		"Preserve the original meaning and tone.\n\n" + text

	req := &api.GenerateRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: new(bool),
	}

	ctx := context.Background()
	rspText := ""
	respFunc := func(resp api.GenerateResponse) error {
		rspText = resp.Response
		return nil
	}

	err := c.Client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}
	return rspText, nil
}
