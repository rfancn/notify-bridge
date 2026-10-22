APP := notify-bridge
MODULE := github.com/rfancn/notify-bridge

VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
  -X $(MODULE)/cmd.version=$(VERSION) \
  -X $(MODULE)/cmd.commit=$(COMMIT) \
  -X $(MODULE)/cmd.buildDate=$(BUILD_DATE)

.PHONY: build

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(APP) .
	@echo "built ./$(APP)"