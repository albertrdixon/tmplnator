PROJECT = github.com/albertrdixon/tmplnator
LDFLAGS = "-X $(PROJECT)/version.Build $$(git rev-parse --short HEAD) -s"
BINARY = "cmd/t2/t2.go"
TEST_COMMAND = TNATOR_DIR=$(shell pwd)/fixtures godep go test

.PHONY: dep-save dep-restore test test-verbose test-integration vet lint build install clean

all: test

help:
	@echo "Available targets:"
	@echo ""
	@echo "  dep-save"
	@echo "  dep-restore"
	@echo "  test"
	@echo "  test-verbose"
	@echo "  test-integration"
	@echo "  vet"
	@echo "  lint"
	@echo "  build"
	@echo "  build-docker"
	@echo "  install"
	@echo "  clean"

dep-save:
	godep save ./...

dep-restore:
	godep restore

test:
	$(TEST_COMMAND) ./...

test-verbose:
	$(TEST_COMMAND) -test.v ./...

test-integration:
	$(TEST_COMMAND) ./... -tags integration

vet:
	go vet ./...

lint:
	golint ./...

build:
	godep go build -ldflags $(LDFLAGS) $(BINARY)

install:
	godep go install -ldflags $(LDFLAGS) $(BINARIES)

clean:
	go clean ./...
	rm -rf t2
