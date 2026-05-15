#!/usr/bin/env bash
# SPDX-FileCopyrightText: © 2023 OneEyeFPV oneeyefpv@gmail.com
# SPDX-License-Identifier: GPL-3.0-or-later
# SPDX-License-Identifier: FS-0.9-or-later
#
# Cross-compiles a Windows amd64 binary from a Linux host using
# Dockerfile.windows-amd64 (Ubuntu + MinGW-w64). No Windows host required.
#
# Usage: ./scripts/build-windows-amd64.sh

set -eu -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOCKERFILE="${ROOT_DIR}/Dockerfile.windows-amd64"
IMAGE_NAME="elrs-joystick-control-build-windows-amd64"
OUT_PATH="elrs-joystick-control-windows-amd64.exe"
CACHE_DIR="${ROOT_DIR}/.cache/windows-amd64"

mkdir -p "${CACHE_DIR}"

echo "Building Docker image ${IMAGE_NAME} using ${DOCKERFILE}..."
if docker buildx version >/dev/null 2>&1; then
  docker buildx build --platform linux/amd64 --load -f "${DOCKERFILE}" -t "${IMAGE_NAME}" "${ROOT_DIR}"
else
  docker build --platform=linux/amd64 -f "${DOCKERFILE}" -t "${IMAGE_NAME}" "${ROOT_DIR}"
fi

echo "Running build inside container..."
docker run --rm \
  -v "${ROOT_DIR}":/src \
  -v "${CACHE_DIR}":/cache \
  -w /src \
  -e GOMODCACHE=/cache/gomod \
  -e GOCACHE=/cache/gocache \
  -e NPM_CONFIG_CACHE=/cache/npm \
  -e XDG_CACHE_HOME=/cache \
  "${IMAGE_NAME}" \
  "GOOS=linux GOARCH=amd64 CC=gcc go generate ./... && go build -tags dynamic --ldflags '-s -w' -o '${OUT_PATH}' ./cmd/elrs-joystick-control"

echo "Build finished: ${OUT_PATH}"
