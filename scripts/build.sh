#!/bin/bash

# ensure local cache dir exists for node/go deps
mkdir -p .cache

docker run --rm -it -m 4096m \
	-v "$(pwd)":/app \
	-v "$(pwd)/.cache":/cache \
	-e GOMODCACHE=/cache/gomod \
	-e GOCACHE=/cache/gocache \
	-e NPM_CONFIG_CACHE=/cache/npm \
	-e XDG_CACHE_HOME=/cache \
	oneeyefpv/linux-amd64-builder \
	"cd /app && go generate ./... && go build -tags static -o elrs-joystick-control ./cmd/elrs-joystick-control"