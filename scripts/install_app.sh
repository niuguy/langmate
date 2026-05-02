#!/bin/bash

# Build LangMate.app and install it into /Applications for local upgrades.

set -euo pipefail

APP_NAME="LangMate"
APP_BUNDLE="$APP_NAME.app"
INSTALL_DIR="/Applications"
INSTALL_PATH="$INSTALL_DIR/$APP_BUNDLE"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR"

echo "Building $APP_BUNDLE..."
"$ROOT_DIR/scripts/build_app.sh"

echo "Quitting running $APP_NAME instances..."
osascript -e "tell application \"$APP_NAME\" to quit" >/dev/null 2>&1 || true
pkill -x langmate >/dev/null 2>&1 || true

sleep 1

if [[ -d "$INSTALL_PATH" ]]; then
  echo "Replacing $INSTALL_PATH..."
  rm -rf "$INSTALL_PATH"
else
  echo "Installing to $INSTALL_PATH..."
fi

ditto "$ROOT_DIR/$APP_BUNDLE" "$INSTALL_PATH"

echo "Launching $APP_NAME..."
open "$INSTALL_PATH"

echo ""
echo "Installed and launched $INSTALL_PATH"
