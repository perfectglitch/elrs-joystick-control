#!/bin/bash

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOCKERFILE="${ROOT_DIR}/Dockerfile.linux-amd64"
IMAGE_NAME="elrs-joystick-control-build-amd64"
OUT_PATH="elrs-joystick-control-linux-amd64"
CACHE_DIR="${ROOT_DIR}/.cache/amd64"

# ensure local cache dir exists for node/go deps
mkdir -p .cache

echo "Building Docker image ${IMAGE_NAME} using ${DOCKERFILE}..."
if docker buildx version >/dev/null 2>&1; then
  docker buildx build --platform linux/amd64 --load -f "${DOCKERFILE}" -t "${IMAGE_NAME}" "${ROOT_DIR}"
else
  docker build --platform=linux/amd64 -f "${DOCKERFILE}" -t "${IMAGE_NAME}" "${ROOT_DIR}" || true
fi

docker run --rm -it -m 4096m \
	-v "$(pwd)":/app \
	-v "${CACHE_DIR}":/cache \
	-e GOMODCACHE=/cache/gomod \
	-e GOCACHE=/cache/gocache \
	-e NPM_CONFIG_CACHE=/cache/npm \
	-e XDG_CACHE_HOME=/cache \
	-e GOPATH=/go \
	-e GOROOT=/go-sdk \
	-e CGO_ENABLED=1 \
	-e GOOS=linux \
	-e GOARCH=amd64 \
	${IMAGE_NAME} \
	"cd /app && go generate ./... && go build -tags static --ldflags '-s -w' -o ${OUT_PATH} ./cmd/elrs-joystick-control"