#!/bin/bash

# Build LangMate.app macOS application bundle

set -e

APP_NAME="LangMate"
VERSION="${VERSION:-1.0.4}"
APP_DIR="$APP_NAME.app"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"

# Clean previous build
rm -rf "$APP_DIR"

# Create directory structure
mkdir -p "$MACOS_DIR"
mkdir -p "$RESOURCES_DIR"

# Build the Go binary
echo "Building langmate binary..."
go build -o "$MACOS_DIR/langmate" .

echo "Building preview helper..."
swiftc "$PWD/preview-helper/main.swift" -o "$MACOS_DIR/langmate-preview"

# Create Info.plist
cat > "$CONTENTS_DIR/Info.plist" << 'EOF'
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
    <string>__VERSION__</string>
    <key>CFBundleShortVersionString</key>
    <string>__VERSION__</string>
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
EOF

sed -i '' "s/__VERSION__/$VERSION/g" "$CONTENTS_DIR/Info.plist"

# No launcher script needed - binary auto-detects .app bundle

echo ""
echo "✓ Built $APP_DIR successfully!"
echo ""
echo "To install:"
echo "  mv $APP_DIR /Applications/"
echo ""
echo "Then:"
echo "  1. Open System Settings > Privacy & Security > Accessibility"
echo "  2. Click + and add LangMate from Applications"
echo "  3. Double-click LangMate.app to start"
echo ""
echo "The app runs in the background. Use Cmd+Ctrl+R to rephrase selected text."
