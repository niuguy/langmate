class Langmate < Formula
  desc "Fast macOS menu bar app for AI-powered selected-text rephrasing"
  homepage "https://github.com/niuguy/langmate"
  url "https://github.com/niuguy/langmate/archive/refs/tags/v1.0.3.tar.gz"
  sha256 "0000000000000000000000000000000000000000000000000000000000000000"
  license "MIT"
  head "https://github.com/niuguy/langmate.git", branch: "main"

  depends_on "go" => :build
  depends_on xcode: ["14.0", :build]

  def install
    app = prefix/"LangMate.app"
    contents = app/"Contents"
    macos = contents/"MacOS"
    resources = contents/"Resources"

    macos.mkpath
    resources.mkpath

    system "go", "build", "-trimpath", "-ldflags", "-s -w", "-o", macos/"langmate", "."
    bin.install macos/"langmate"

    system "swiftc", "preview-helper/main.swift", "-o", macos/"langmate-preview"

    resources.install "app/icon.png"

    (contents/"Info.plist").write <<~XML
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
      <dict>
          <key>CFBundleName</key>
          <string>LangMate</string>
          <key>CFBundleDisplayName</key>
          <string>LangMate</string>
          <key>CFBundleIdentifier</key>
          <string>com.langmate.app</string>
          <key>CFBundleVersion</key>
          <string>1.0.3</string>
          <key>CFBundleShortVersionString</key>
          <string>1.0.3</string>
          <key>CFBundlePackageType</key>
          <string>APPL</string>
          <key>CFBundleExecutable</key>
          <string>langmate</string>
          <key>LSMinimumSystemVersion</key>
          <string>10.15</string>
          <key>LSUIElement</key>
          <true/>
          <key>NSHighResolutionCapable</key>
          <true/>
          <key>LSApplicationCategoryType</key>
          <string>public.app-category.productivity</string>
      </dict>
      </plist>
    XML

    installer = buildpath/"scripts/langmate-install-app"
    inreplace installer, "__PREFIX__", prefix
    bin.install installer
  end

  def caveats
    <<~EOS
      To install the app bundle into /Applications and launch it:
        langmate-install-app

      Then grant Accessibility permission:
        System Settings -> Privacy & Security -> Accessibility -> LangMate

      Configure OpenAI:
        echo 'OPENAI_API_KEY=sk-your-api-key-here' > ~/.langmate.env

      For local Ollama presets, pull a model first:
        ollama pull qwen3:8b
    EOS
  end

  test do
    assert_path_exists prefix/"LangMate.app/Contents/MacOS/langmate"
    assert_path_exists prefix/"LangMate.app/Contents/MacOS/langmate-preview"
    assert_match "langmate", shell_output("#{bin}/langmate --help")
  end
end
