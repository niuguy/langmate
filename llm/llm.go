package llm

type TextProcessor interface {
	TransferText(text string, lang string) (string, error)
}

func CreateTextProcessor(llmtype string) TextProcessor {
	switch llmtype {
	case "o", "openai":
		return NewOpenAIClient()
	case "l", "ollama":
		return NewOllamaClient()
	default:
		return NewOpenAIClient()
	}
}
