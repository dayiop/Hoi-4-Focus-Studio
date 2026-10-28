#!/usr/bin/env bash
# Builds portable archives and Debian packages for Linux x86_64 / ARM64.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/dist"
mkdir -p "$OUT"
for ARCH in amd64 arm64; do
  (
    cd "$ROOT/desktop-src"
    GOOS=linux GOARCH="$ARCH" CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o "$OUT/HOI4-Focus-Studio-linux-$ARCH" .
  )
  STAGE="$OUT/stage-$ARCH"
  rm -rf "$STAGE"
  mkdir -p "$STAGE/opt/hoi4-focus-studio" "$STAGE/usr/share/applications" "$STAGE/usr/share/icons/hicolor/256x256/apps" "$STAGE/usr/bin" "$STAGE/DEBIAN"
  cp "$OUT/HOI4-Focus-Studio-linux-$ARCH" "$STAGE/opt/hoi4-focus-studio/HOI4-Focus-Studio"
  chmod 755 "$STAGE/opt/hoi4-focus-studio/HOI4-Focus-Studio"
  ln -s /opt/hoi4-focus-studio/HOI4-Focus-Studio "$STAGE/usr/bin/hoi4-focus-studio"
  cp "$ROOT/desktop-src/assets/hoi4-focus-studio.png" "$STAGE/usr/share/icons/hicolor/256x256/apps/hoi4-focus-studio.png"
  cp "$ROOT/desktop-src/linux/hoi4-focus-studio.desktop" "$STAGE/usr/share/applications/"
  cat > "$STAGE/DEBIAN/control" <<EOF
Package: hoi4-focus-studio
Version: 1.1.0
Section: graphics
Priority: optional
Architecture: $ARCH
Maintainer: HOI4 Focus Studio Community <noreply@example.invalid>
Description: Visual national focus tree editor for Hearts of Iron IV
 Create and edit focus trees, auto-save projects, export HOI4 scripts,
 and transfer projects with JSON files. Uses a locally installed browser.
Depends: xdg-utils
EOF
  chmod 644 "$STAGE/DEBIAN/control"
  chmod g-s "$STAGE/DEBIAN"
  dpkg-deb --root-owner-group --build "$STAGE" "$OUT/HOI4-Focus-Studio-linux-$ARCH.deb"
  PORTABLE="$OUT/HOI4-Focus-Studio-linux-$ARCH-portable"
  rm -rf "$PORTABLE"
  mkdir -p "$PORTABLE"
  cp "$OUT/HOI4-Focus-Studio-linux-$ARCH" "$PORTABLE/HOI4-Focus-Studio"
  cp "$ROOT/desktop-src/assets/hoi4-focus-studio.png" "$PORTABLE/icon.png"
  cp "$ROOT/desktop-src/linux/README.txt" "$PORTABLE/README.txt"
  tar -C "$OUT" -czf "$OUT/HOI4-Focus-Studio-linux-$ARCH-portable.tar.gz" "$(basename "$PORTABLE")"
  rm -rf "$STAGE" "$PORTABLE"
done
echo 'Created Linux .deb and portable archives for amd64 and arm64'
