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

	prompt := fmt.Sprintf("You will receive a text and a destination language. "+
		"If the text is not in the destination language, translate it."+
		"If the text is already in the destination language, rephrase it for clarity and style."+
		"If the text is a single word or slang, please explain it in a way dictionary does,"+
		"When the input is word or slang, in the bottom please list some alternatives if there are any"+
		"\n\nInput Text: \" %s \"\nDestination Language: \" %s \"\n\nOutput:", text, lang)

	req := &api.GenerateRequest{
		Model:  "llama3",
		Prompt: prompt,

		// set streaming to false
		Stream: new(bool),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		fmt.Println(resp.Response)
		return nil
	}

	err := c.Client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}
	return "", nil

}
