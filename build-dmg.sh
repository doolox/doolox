#!/bin/sh
# Build a macOS .app bundle and a .dmg disk image for Doolox.
#
# Requires: a Doolox binary at bin/doolox (native or universal) and macOS
# (hdiutil). The admin web assets are bundled into the .app so the dashboard
# renders its CSS/JS.
#
# Usage:    ./build-dmg.sh
# Output:   dist/Doolox-<version>.dmg (contains Doolox.app + /Applications link)
set -eu

ROOT="$(cd "$(dirname "$0")" && pwd)"
APP_NAME="Doolox"
VERSION="$(sed -n 's/^#define MyAppVersion "\([^"]*\)".*/\1/p' "$ROOT/version.iss" | head -n1)"

if [ -z "$VERSION" ]; then
    echo "error: could not read MyAppVersion from $ROOT/version.iss" >&2
    exit 1
fi

BINARY="$ROOT/bin/doolox"
if [ ! -f "$BINARY" ]; then
    echo "error: $BINARY not found — build the macOS binary first (make build)" >&2
    exit 1
fi

STAGE="$ROOT/build/macos"
APP_DIR="$STAGE/$APP_NAME.app"
CONTENTS="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS/MacOS"
RESOURCES="$CONTENTS/Resources"
OUT_DIR="$ROOT/dist"
DMG="$OUT_DIR/${APP_NAME}-${VERSION}.dmg"

rm -rf "$STAGE"
mkdir -p "$MACOS_DIR" "$RESOURCES"

# 1. The executable and the admin web assets (served next to the binary).
install -m 0755 "$BINARY" "$MACOS_DIR/doolox"
mkdir -p "$MACOS_DIR/admin"
cp -R "$ROOT/admin/theme" "$MACOS_DIR/admin/theme"

# 2. Info.plist — LSUIElement hides the Dock icon (Doolox is a menu-bar app).
cat > "$CONTENTS/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleDisplayName</key>
    <string>$APP_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>com.doolox.app</string>
    <key>CFBundleVersion</key>
    <string>$VERSION</string>
    <key>CFBundleShortVersionString</key>
    <string>$VERSION</string>
    <key>CFBundleExecutable</key>
    <string>doolox</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
EOF

# 3. Assemble the disk image with the .app and a /Applications symlink.
STAGING="$STAGE/dmg"
mkdir -p "$STAGING"
cp -R "$APP_DIR" "$STAGING/"
ln -s /Applications "$STAGING/Applications"

mkdir -p "$OUT_DIR"
rm -f "$DMG"
hdiutil create -volname "$APP_NAME" -srcfolder "$STAGING" -ov -format UDZO "$DMG"

echo "Built $DMG"
