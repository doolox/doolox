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

# 1. The executable (renamed doolox-bin; a launcher wrapper seeds the bundled
#    defaults on first run, then execs it) and the admin web assets.
install -m 0755 "$BINARY" "$MACOS_DIR/doolox-bin"
mkdir -p "$MACOS_DIR/admin"
cp -R "$ROOT/admin/theme" "$MACOS_DIR/admin/theme"

# Seed data (defaults/) used to populate ~/doolox on first run.
mkdir -p "$RESOURCES/seed"
cp -R "$ROOT/defaults/simple/." "$RESOURCES/seed/simple/"
cp -R "$ROOT/defaults/append/." "$RESOURCES/seed/append/"
cp "$ROOT/defaults/config.json" "$RESOURCES/seed/config.json"

# Launcher wrapper (CFBundleExecutable) that seeds the bundled defaults on
# first run, then execs the real binary. It mirrors the Linux launcher.
cat > "$MACOS_DIR/doolox" <<'EOF'
#!/bin/sh
DIR="$(cd "$(dirname "$0")" && pwd)"
DATA_DIR="${DOOLOX_HOME:-$HOME/doolox}"
SITE_DIR="$DATA_DIR/websites"
mkdir -p "$DATA_DIR" "$DATA_DIR/templates"
for SITE in simple append; do
    DEST="$SITE_DIR/$SITE"
    mkdir -p "$DEST"
    if [ ! -d "$DEST/src" ] && [ ! -d "$DEST/content" ]; then
        cp -R "$DIR/../Resources/seed/$SITE/." "$DEST/"
    fi
done
if [ ! -f "$DATA_DIR/config.json" ]; then
    cp "$DIR/../Resources/seed/config.json" "$DATA_DIR/"
fi
exec "$DIR/doolox-bin" "$@"
EOF
chmod 0755 "$MACOS_DIR/doolox"

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
