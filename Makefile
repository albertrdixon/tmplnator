PROJECT = github.com/albertrdixon/tmplnator
INSTALL = $(PROJECT)/cmd/t2
EXECUTABLE = "t2"
LDFLAGS = "-s"
BINARY = "cmd/t2/t2.go"
TEST_COMMAND = godep go test
PLATFORM = "$$(echo "$$(uname)" | tr '[A-Z]' '[a-z]')"
VERSION = "$$(./t2 -v)"
BUILD_ARGS = ""

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
	@echo "==> Running all tests"
	@echo ""
	@$(TEST_COMMAND) ./...

test-verbose:
	@echo "==> Running all tests (verbose output)"
	@ echo ""
	@$(TEST_COMMAND) -test.v ./...

test-integration:
	$(TEST_COMMAND) ./... -tags integration

vet:
	go vet ./...

lint:
	golint ./...

build:
	@echo "==> Building $(EXECUTABLE) with ldflags '$(LDFLAGS)'"
	@godep go build -ldflags $(LDFLAGS) $(BINARY)

install:
	@echo "==> Installing $(EXECUTABLE) with ldflags $(LDFLAGS)"
	@godep go install -ldflags $(LDFLAGS) $(INSTALL)

package: build
	@echo "==> Tar'ing up the binary"
	@test -f t2 && tar czf tnator-$(PLATFORM)-amd64.tar.gz t2

clean:
	go clean ./...
	rm -rf t2
