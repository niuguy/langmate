package llm

type TextProcessor interface {
	TransferText(text string, lang string) (string, error)
}

func CreateTextProcessor(llmtype string) TextProcessor {
	switch llmtype {
	case "gpt":
		return NewOpenAIClient()
	case "llama":
		return NewOllamaClient()
	default:
		return NewOpenAIClient()
	}
}
