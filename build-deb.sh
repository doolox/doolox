#!/bin/sh
# Build a .deb package for Doolox (Debian / Ubuntu / Linux Mint).
#
# Requires: go (with a C compiler for cgo/SQLite) and dpkg-deb.
# Usage:    ./build-deb.sh
# Output:   dist/doolox_1.0.0_amd64.deb
set -eu

VERSION="1.0.0"
ARCH="amd64"
PKG="doolox"

ROOT="$(cd "$(dirname "$0")" && pwd)"
STAGE="$ROOT/build/deb"
OUT_DIR="$ROOT/dist"
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

# 4. Read-only seed data used to populate ~/Websites on first run.
mkdir -p "$STAGE/opt/doolox/seed"
cp -R "$ROOT/content" "$STAGE/opt/doolox/seed/content"
cp "$ROOT/config.json" "$STAGE/opt/doolox/seed/config.json"
cp "$ROOT/doolox.db" "$STAGE/opt/doolox/seed/doolox.db"

# 5. Icons -> hicolor icon theme (referenced by the .desktop file as "doolox").
mkdir -p "$STAGE/usr/share/icons/hicolor/scalable/apps"
mkdir -p "$STAGE/usr/share/icons/hicolor/32x32/apps"
cp "$ROOT/assets/trayicon.svg" "$STAGE/usr/share/icons/hicolor/scalable/apps/doolox.svg"
cp "$ROOT/assets/trayicon.png" "$STAGE/usr/share/icons/hicolor/32x32/apps/doolox.png"

# 6. Make the launcher and maintainer scripts executable.
chmod 0755 "$STAGE/usr/bin/doolox"
chmod 0755 "$STAGE/DEBIAN/postinst" "$STAGE/DEBIAN/postrm"

# 7. Fill in Installed-Size (in KiB) so dpkg-deb doesn't warn.
SIZE_KB="$(du -sk "$STAGE" | cut -f1)"
sed -i '/^Installed-Size:/d' "$STAGE/DEBIAN/control"
printf 'Installed-Size: %s\n' "$SIZE_KB" >> "$STAGE/DEBIAN/control"

# 8. Build the package.
mkdir -p "$OUT_DIR"
dpkg-deb --build "$STAGE" "$OUT_DIR/$DEB_NAME"

echo "Built $OUT_DIR/$DEB_NAME"
