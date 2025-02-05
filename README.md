# LangMate

LangMate is a command-line application developed in Go, designed to provide users with a seamless interaction with different language models. It supports various models, allowing users to choose between models like Ollama and OpenAI for processing and generating text.




https://github.com/user-attachments/assets/f208a9f9-70b6-482b-adc4-127e9a9ef226







## Features

- **Multiple Language Models**: Choose between  OpenAI and Ollama 
- **Easy Text Processing**: Simply perform a double `Command+C` on any text to translate or rephrase it.


## Installation

### Prerequisites

- Go 1.16 or higher

### Installing from Source

To install LangMate from source, follow these steps:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/niuguy/langmate.git
   cd langmate
   ```

2. **Build the application:**

   ```bash
   go build -o langmate
   ```

3. **Optionally, install the application globally:**

   ```bash
   go install
   ```

### Direct Installation

If you prefer not to clone the repository, you can install directly using Go:

```bash
go install github.com/niuguy/langmate@latest
```

Ensure your `GOPATH/bin` is in your system's PATH to run the application from any terminal.

## Usage

To use LangMate, run the executable with the desired options:

```bash
langmate [-m model] [--lang language]
```

### Options

- `-m, --model`: Specify the model to use (default: "gpt"). Available models include "gpt", "llama" , gpt represents OpenAI's GPT-4 Turbo model, and llama represents Ollama's llama3-8b.
- `-l, --lang`: Specify the target language (default: "en"). Available languages include "en", "fr", etc.

### Examples

- **Using the default model (GPT-4 Turbo) and language (English):**

  ```bash
  langmate
  ```

- **Using a different model and language:**

  ```bash
  langmate -m gpt --lang fr
  ```


## Configuration

Set environment variables as needed:

```bash
export OPENAI_API_KEY="your-openai-api-key"
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have feedback or proposals for new features.

## License

Distributed under the MIT License. See `LICENSE` for more information.
