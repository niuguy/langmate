package llm

import "context"

type TextProcessor interface {
	TransferText(ctx context.Context, text string, lang string) (string, error)
	RephraseText(ctx context.Context, text string, lang string) (string, error)
}

type ModelPreset struct {
	ID       string
	Title    string
	Provider string
	Model    string
}

const (
	ProviderOpenAI = "openai"
	ProviderOllama = "ollama"
)

var ModelPresets = []ModelPreset{
	{
		ID:       "openai-fast",
		Title:    "OpenAI Fast (GPT-5 Nano)",
		Provider: ProviderOpenAI,
		Model:    "gpt-5-nano",
	},
	{
		ID:       "openai-balanced",
		Title:    "OpenAI Balanced (GPT-4.1 Mini)",
		Provider: ProviderOpenAI,
		Model:    "gpt-4.1-mini",
	},
	{
		ID:       "openai-quality",
		Title:    "OpenAI Quality (GPT-5 Mini)",
		Provider: ProviderOpenAI,
		Model:    "gpt-5-mini",
	},
	{
		ID:       "ollama-local",
		Title:    "Ollama Local (Qwen3 8B)",
		Provider: ProviderOllama,
		Model:    "qwen3:8b",
	},
	{
		ID:       "ollama-light",
		Title:    "Ollama Light (Qwen3 4B)",
		Provider: ProviderOllama,
		Model:    "qwen3:4b",
	},
	{
		ID:       "ollama-legacy",
		Title:    "Ollama Legacy (Llama 3)",
		Provider: ProviderOllama,
		Model:    "llama3",
	},
}

func NormalizeModelPresetID(id string) string {
	switch id {
	case "", "gpt":
		return "openai-balanced"
	case "llama":
		return "ollama-local"
	default:
		if _, ok := FindModelPreset(id); ok {
			return id
		}
		return "openai-balanced"
	}
}

func FindModelPreset(id string) (ModelPreset, bool) {
	for _, preset := range ModelPresets {
		if preset.ID == id {
			return preset, true
		}
	}
	return ModelPreset{}, false
}

func CreateTextProcessor(modelID string) (TextProcessor, error) {
	preset, _ := FindModelPreset(NormalizeModelPresetID(modelID))
	switch preset.Provider {
	case ProviderOpenAI:
		return NewOpenAIClient(preset.Model)
	case ProviderOllama:
		return NewOllamaClient(preset.Model)
	default:
		return NewOpenAIClient("gpt-4.1-mini")
	}
}
