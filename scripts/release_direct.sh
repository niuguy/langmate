#!/bin/bash

# Build, sign, package, and optionally notarize LangMate for direct macOS distribution.

set -euo pipefail

APP_NAME="${APP_NAME:-LangMate}"
BUNDLE_ID="${BUNDLE_ID:-com.langmate.app}"
VERSION="${VERSION:-1.0.4}"
BUILD="${BUILD:-$VERSION}"
MIN_MACOS_VERSION="${MIN_MACOS_VERSION:-10.15}"
SIGN_IDENTITY="${SIGN_IDENTITY:?Set SIGN_IDENTITY to your Developer ID Application certificate name}"
NOTARY_PROFILE="${NOTARY_PROFILE:-}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="$ROOT_DIR/dist"
APP_DIR="$DIST_DIR/$APP_NAME.app"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"
ENTITLEMENTS="$ROOT_DIR/scripts/direct_distribution.entitlements"
DMG_PATH="$DIST_DIR/$APP_NAME-$VERSION.dmg"
STAGING_DMG="$DIST_DIR/$APP_NAME-$VERSION-staging.dmg"

rm -rf "$DIST_DIR"
mkdir -p "$MACOS_DIR" "$RESOURCES_DIR"

echo "Building $APP_NAME..."
go build -trimpath -ldflags "-s -w" -o "$MACOS_DIR/langmate" "$ROOT_DIR"

echo "Building preview helper..."
swiftc "$ROOT_DIR/preview-helper/main.swift" -o "$MACOS_DIR/langmate-preview"

cp "$ROOT_DIR/app/icon.png" "$RESOURCES_DIR/icon.png"

cat > "$CONTENTS_DIR/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleDisplayName</key>
    <string>$APP_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>$BUNDLE_ID</string>
    <key>CFBundleVersion</key>
    <string>$BUILD</string>
    <key>CFBundleShortVersionString</key>
    <string>$VERSION</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleExecutable</key>
    <string>langmate</string>
    <key>LSMinimumSystemVersion</key>
    <string>$MIN_MACOS_VERSION</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSApplicationCategoryType</key>
    <string>public.app-category.productivity</string>
    <key>NSAppleEventsUsageDescription</key>
    <string>LangMate sends copy and paste commands to the active app only when you press its rephrase hotkey.</string>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright (c) 2026 LangMate. All rights reserved.</string>
</dict>
</plist>
EOF

echo "Signing executable..."
codesign --force \
  --options runtime \
  --timestamp \
  --entitlements "$ENTITLEMENTS" \
  --sign "$SIGN_IDENTITY" \
  "$MACOS_DIR/langmate"

echo "Signing preview helper..."
codesign --force \
  --options runtime \
  --timestamp \
  --sign "$SIGN_IDENTITY" \
  "$MACOS_DIR/langmate-preview"

echo "Signing app bundle..."
codesign --force \
  --options runtime \
  --timestamp \
  --entitlements "$ENTITLEMENTS" \
  --sign "$SIGN_IDENTITY" \
  "$APP_DIR"

echo "Verifying signature..."
codesign --verify --deep --strict --verbose=2 "$APP_DIR"
spctl --assess --type execute --verbose=4 "$APP_DIR" || true

echo "Creating DMG..."
hdiutil create \
  -volname "$APP_NAME" \
  -srcfolder "$APP_DIR" \
  -ov \
  -format UDRW \
  "$STAGING_DMG"

hdiutil convert "$STAGING_DMG" \
  -format UDZO \
  -imagekey zlib-level=9 \
  -o "$DMG_PATH"

rm -f "$STAGING_DMG"

echo "Signing DMG..."
codesign --force \
  --timestamp \
  --sign "$SIGN_IDENTITY" \
  "$DMG_PATH"

echo "Verifying DMG signature..."
codesign --verify --verbose=2 "$DMG_PATH"

if [[ -n "$NOTARY_PROFILE" ]]; then
  echo "Submitting DMG for notarization..."
  xcrun notarytool submit "$DMG_PATH" \
    --keychain-profile "$NOTARY_PROFILE" \
    --wait

  echo "Stapling notarization ticket..."
  xcrun stapler staple "$DMG_PATH"
  xcrun stapler validate "$DMG_PATH"
else
  echo "Skipping notarization because NOTARY_PROFILE is not set."
fi

echo ""
echo "Release artifact:"
echo "  $DMG_PATH"
