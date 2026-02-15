#!/usr/bin/env bash
set -e

APP_NAME="flux"

if [ -z "$1" ]; then
  echo "Usage: ./build.sh vX.Y.Z"
  exit 1
fi

VERSION=$1
BUILD_DIR="dist/$VERSION"

echo "Building $APP_NAME version: $VERSION"
echo "Output directory: $BUILD_DIR"

rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS="
-s -w
-X main.version=$VERSION
-X main.commit=$COMMIT
-X main.buildDate=$BUILD_DATE
"

platforms=(
  "linux/amd64"
  "linux/arm64"
  "linux/arm"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

for platform in "${platforms[@]}"
do
  GOOS=${platform%/*}
  GOARCH=${platform#*/}

  OUTPUT_NAME="$APP_NAME-$GOOS-$GOARCH"

  if [ "$GOOS" = "windows" ]; then
    OUTPUT_NAME+=".exe"
  fi

  echo "→ Building $GOOS/$GOARCH"

  GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
    go build -ldflags="$LDFLAGS" \
    -o "$BUILD_DIR/$OUTPUT_NAME" \
    ./cmd/agent
done

echo "Generating checksums..."
cd "$BUILD_DIR"
sha256sum * > checksums.txt
cd - > /dev/null

echo ""
echo "Build complete."
echo "Artifacts in: $BUILD_DIR"