#!/usr/bin/env bash
set -eu -o pipefail

# Builds an arm64 (aarch64) linux binary using Docker and Dockerfile.linux-arm64
# Usage: ./scripts/build-arm64.sh [output-path]

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOCKERFILE="${ROOT_DIR}/Dockerfile.linux-arm64"
IMAGE_NAME="elrs-joystick-control-build-arm64"
OUT_PATH="elrs-joystick-control-linux-arm64"

mkdir -p .cache

echo "Building Docker image ${IMAGE_NAME} using ${DOCKERFILE}..."
if docker buildx version >/dev/null 2>&1; then
  docker buildx build --platform linux/arm64 --load -f "${DOCKERFILE}" -t "${IMAGE_NAME}" "${ROOT_DIR}"
else
  docker build --platform=linux/arm64 -f "${DOCKERFILE}" -t "${IMAGE_NAME}" "${ROOT_DIR}" || true
fi

echo "Running build inside container..."
docker run --rm \
  -v "${ROOT_DIR}":/src \
  -v "$(pwd)/.cache":/cache \
  -w /src \
  -e GOMODCACHE=/cache/gomod \
  -e GOCACHE=/cache/gocache \
  -e NPM_CONFIG_CACHE=/cache/npm \
  -e XDG_CACHE_HOME=/cache \
  -e GOPATH=/go \
  -e GOROOT=/go-sdk \
  -e CGO_ENABLED=1 \
  -e GOOS=linux \
  -e GOARCH=arm64 \
  "${IMAGE_NAME}" "cd /src && go generate ./...  && /go-sdk/bin/go build -v -o '${OUT_PATH}' ./cmd/elrs-joystick-control"

echo "Build finished: ${OUT_PATH}"
