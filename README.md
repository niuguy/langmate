# LangMate

LangMate is a command-line application developed in Go, designed to provide users with a seamless interaction with different language models. It supports various models, allowing users to choose between models like Ollama and OpenAI for processing and generating text.

![image](https://github.com/user-attachments/assets/57d1df13-7968-49dc-afef-4dca71f4d1d9)


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

To use LangMate, run the executable followed by the options:

```bash
langmate [options]
```
by default, the application uses the OpenAI model to process text. To use the local Ollama model, start your ollama server then run with

```bash
langmate l
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
