# LangMate

LangMate is a macOS app that rephrases selected text in place using AI. Select any text, press **Cmd+Ctrl+R**, and watch it get rephrased instantly.


## Features

- **Rephrase in Place**: Select text in any app, press Cmd+Ctrl+R, and the text is replaced with a rephrased version
- **Menu Bar App**: Runs quietly in the background with a menu bar indicator
- **Multiple Language Models**: Supports OpenAI GPT-4.1 and Ollama (llama3)

## Installation

### Prerequisites

- macOS 10.15 or higher
- Go 1.22 or higher
- OpenAI API key

### Homebrew

Install the app:

```bash
brew tap niuguy/langmate
brew install --cask langmate
```

For source builds instead of the cask:

```bash
brew install niuguy/langmate/langmate
langmate-install-app
```

### Build from Source

1. **Clone the repository:**

   ```bash
   git clone https://github.com/niuguy/langmate.git
   cd langmate
   ```

2. **Build the macOS app:**

   ```bash
   ./scripts/build_app.sh
   ```

3. **Install to Applications:**

   ```bash
   cp -r LangMate.app /Applications/
   ```

4. **Configure your API key:**

   ```bash
   echo 'OPENAI_API_KEY=sk-your-api-key-here' > ~/.langmate.env
   ```

5. **Grant Accessibility permission:**
   - Open **System Settings** → **Privacy & Security** → **Accessibility**
   - Click **+** and add `/Applications/LangMate.app`
   - Enable the toggle

6. **Launch the app** from Spotlight or Applications folder

### Command Line Usage

You can also run LangMate directly from the terminal:

```bash
# Build the binary
go build -o langmate

# Run with default settings (GPT-4.1, English)
./langmate

# Run with different model or language
./langmate -m llama -l fr
```

## Usage

1. Open LangMate (it appears as "LM" in your menu bar)
2. Select any text in any application
3. Press **Cmd+Ctrl+R**
4. The menu bar shows "Rephrasing..." while processing
5. Your selected text is replaced with the rephrased version
6. Your original clipboard contents are restored after the paste completes

## Options

- `-m, --model`: Model to use - `gpt` (OpenAI GPT-4.1) or `llama` (Ollama llama3). Default: `gpt`
- `-l, --lang`: Target language code (e.g., `en`, `fr`, `es`). Default: `en`

## Configuration

Create `~/.langmate.env` with your API key:

```
OPENAI_API_KEY=sk-your-api-key-here
```

Alternatively, set the environment variable:

```bash
export OPENAI_API_KEY="sk-your-api-key-here"
```

## Direct Distribution

LangMate can be distributed outside the Mac App Store as a signed and notarized DMG.

Prerequisites:

- Apple Developer Program membership
- A `Developer ID Application` certificate installed in Keychain
- Xcode command line tools

Create a notarytool keychain profile once:

```bash
xcrun notarytool store-credentials langmate-notary \
  --apple-id "you@example.com" \
  --team-id "TEAMID1234" \
  --password "app-specific-password"
```

Build a signed, notarized DMG:

```bash
SIGN_IDENTITY="Developer ID Application: Your Name (TEAMID1234)" \
NOTARY_PROFILE="langmate-notary" \
VERSION="1.0.4" \
./scripts/release_direct.sh
```

The release artifact is written to `dist/LangMate-<version>.dmg`.

LangMate requires Accessibility permission to send the hotkey-driven copy and paste commands. Users should install the app in `/Applications`, launch it once, then enable it in **System Settings** -> **Privacy & Security** -> **Accessibility**.

## Troubleshooting

### Hotkey not working
- Ensure LangMate has Accessibility permission in System Settings
- After rebuilding, remove and re-add the app in Accessibility settings

### App not launching from Spotlight
- Make sure `~/.langmate.env` contains your API key
- Check Console.app for any error messages

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

Distributed under the MIT License. See `LICENSE` for more information.
