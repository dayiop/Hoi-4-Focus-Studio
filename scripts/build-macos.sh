#!/usr/bin/env bash
# Builds a dual-architecture macOS .app, packed for distribution.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/dist"
APP="$OUT/HOI4 Focus Studio.app"
BIN="$APP/Contents/MacOS"
RES="$APP/Contents/Resources"
mkdir -p "$BIN" "$RES"
(
 cd "$ROOT/desktop-src"
 GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o "$BIN/focus-studio-intel" .
 GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o "$BIN/focus-studio-apple-silicon" .
)
cat > "$BIN/HOI4FocusStudio" <<'EOF'
#!/bin/sh
HERE="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd -P)"
case "$(uname -m)" in
  arm64) exec "$HERE/focus-studio-apple-silicon" ;;
  x86_64) exec "$HERE/focus-studio-intel" ;;
  *) echo 'HOI4 Focus Studio: unsupported Mac architecture' >&2; exit 1 ;;
esac
EOF
chmod +x "$BIN/HOI4FocusStudio" "$BIN/focus-studio-intel" "$BIN/focus-studio-apple-silicon"
cp "$ROOT/desktop-src/assets/HOI4FocusStudio.icns" "$RES/HOI4FocusStudio.icns"
cat > "$APP/Contents/Info.plist" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
 <key>CFBundleName</key><string>HOI4 Focus Studio</string>
 <key>CFBundleDisplayName</key><string>HOI4 Focus Studio</string>
 <key>CFBundleIdentifier</key><string>io.github.hoi4focusstudio.desktop</string>
 <key>CFBundleExecutable</key><string>HOI4FocusStudio</string>
 <key>CFBundleIconFile</key><string>HOI4FocusStudio</string>
 <key>CFBundlePackageType</key><string>APPL</string>
 <key>CFBundleShortVersionString</key><string>1.1.0</string>
 <key>CFBundleVersion</key><string>110</string>
 <key>LSMinimumSystemVersion</key><string>11.0</string>
 <key>LSUIElement</key><true/>
 <key>NSHighResolutionCapable</key><true/>
</dict></plist>
EOF
cat > "$RES/INSTALL.txt" <<'EOF'
HOI4 Focus Studio (macOS 11+, Intel / Apple Silicon)

Drag HOI4 Focus Studio.app into /Applications and open it.
Google Chrome, Microsoft Edge, Chromium or Brave: standalone editor window.
Safari/default browser: editor opens in a regular browser tab.
Unsigned app: macOS can request confirmation via Privacy & Security > Open Anyway.
Projects are stored locally in your browser profile. Export JSON backups regularly.
EOF
if command -v codesign >/dev/null 2>&1; then
 codesign --force --deep --sign - "$APP"
fi
rm -f "$OUT/HOI4-Focus-Studio-macOS-universal.zip"
(
 cd "$OUT"
 if command -v ditto >/dev/null 2>&1; then
  ditto -c -k --sequesterRsrc --keepParent 'HOI4 Focus Studio.app' 'HOI4-Focus-Studio-macOS-universal.zip'
 else
  zip -q -r 'HOI4-Focus-Studio-macOS-universal.zip' 'HOI4 Focus Studio.app'
 fi
)
echo "Created $OUT/HOI4-Focus-Studio-macOS-universal.zip"
