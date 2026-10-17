#!/bin/sh
# Build a .deb package for Doolox (Debian / Ubuntu / Linux Mint).
#
# Requires: go (with a C compiler for cgo/SQLite) and dpkg-deb.
# Usage:    ./build-deb.sh
# Output:   dist/doolox_<version>_amd64.deb
set -eu

ROOT="$(cd "$(dirname "$0")" && pwd)"
STAGE="$ROOT/build/deb"
OUT_DIR="$ROOT/dist"

# Version lives in version.iss (single source of truth, shared with installer.iss).
VERSION="$(sed -n 's/^#define MyAppVersion "\([^"]*\)".*/\1/p' "$ROOT/version.iss" | head -n1)"
ARCH="amd64"
PKG="doolox"

if [ -z "$VERSION" ]; then
    echo "error: could not read MyAppVersion from $ROOT/version.iss" >&2
    exit 1
fi

DEB_NAME="${PKG}_${VERSION}_${ARCH}.deb"

# 1. Build the Linux binary (links SQLite via cgo, so a C compiler is needed).
( cd "$ROOT" && make build )

# 2. Assemble the staging directory from the committed packaging files
#    (DEBIAN/ control + maintainer scripts, usr/bin launcher, desktop entry).
rm -rf "$STAGE"
mkdir -p "$STAGE"
cp -R "$ROOT/packaging/debian/." "$STAGE/"

# 3. Application files -> /opt/doolox (the Linux equivalent of Program Files).
mkdir -p "$STAGE/opt/doolox"
install -m 0755 "$ROOT/bin/doolox" "$STAGE/opt/doolox/doolox"
cp -R "$ROOT/assets" "$STAGE/opt/doolox/assets"
mkdir -p "$STAGE/opt/doolox/admin"
cp -R "$ROOT/admin/theme" "$STAGE/opt/doolox/admin/theme"

# 4. Read-only seed data (defaults/) used to populate ~/doolox on first run.
mkdir -p "$STAGE/opt/doolox/seed"
cp -R "$ROOT/defaults/doolox/content" "$STAGE/opt/doolox/seed/content"
cp "$ROOT/defaults/config.json" "$STAGE/opt/doolox/seed/config.json"

# 5. Icons -> hicolor theme + pixmaps fallback (the .desktop uses Icon=doolox).
#    Render PNGs at multiple sizes from the SVG (reusing the project's Go
#    deps, so no ImageMagick/rsvg is needed), then install them along with a
#    pixmaps copy that desktop menus can find even without an icon-theme cache.
( cd "$ROOT" && go run ./tools/rendericons assets/trayicon.svg "$STAGE/.icons" )

for S in 16 24 32 48 64 128 256; do
    mkdir -p "$STAGE/usr/share/icons/hicolor/${S}x${S}/apps"
    cp "$STAGE/.icons/doolox-${S}.png" "$STAGE/usr/share/icons/hicolor/${S}x${S}/apps/doolox.png"
done
mkdir -p "$STAGE/usr/share/icons/hicolor/scalable/apps"
cp "$ROOT/assets/trayicon.svg" "$STAGE/usr/share/icons/hicolor/scalable/apps/doolox.svg"
mkdir -p "$STAGE/usr/share/pixmaps"
cp "$STAGE/.icons/doolox-256.png" "$STAGE/usr/share/pixmaps/doolox.png"
rm -rf "$STAGE/.icons"

# 6. Make the launcher and maintainer scripts executable.
chmod 0755 "$STAGE/usr/bin/doolox"
chmod 0755 "$STAGE/DEBIAN/postinst" "$STAGE/DEBIAN/postrm"

# 7. Substitute the version placeholder in the control file, then fill in
#    Installed-Size (in KiB) so dpkg-deb doesn't warn.
sed -i "s/@VERSION@/$VERSION/" "$STAGE/DEBIAN/control"
SIZE_KB="$(du -sk "$STAGE" | cut -f1)"
sed -i '/^Installed-Size:/d' "$STAGE/DEBIAN/control"
printf 'Installed-Size: %s\n' "$SIZE_KB" >> "$STAGE/DEBIAN/control"

# 8. Build the package.
mkdir -p "$OUT_DIR"
dpkg-deb --build "$STAGE" "$OUT_DIR/$DEB_NAME"

echo "Built $OUT_DIR/$DEB_NAME"
