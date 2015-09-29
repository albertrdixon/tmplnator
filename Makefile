PROJECT = github.com/albertrdixon/tmplnator
TEST_COMMAND = godep go test
EXECUTABLE = t2
PKG = cmd/t2/t2.go
LDFLAGS = -s
PLATFORMS = linux darwin
BUILD_ARGS = ""

.PHONY: dep-save dep-restore test test-verbose build install clean

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
	@echo "==> Saving dependencies..."
	@godep save ./...

dep-restore:
	@echo "==> Restoring dependencies..."
	@godep restore

test:
	@echo "==> Running all tests"
	@echo ""
	@$(TEST_COMMAND) ./...

test-verbose:
	@echo "==> Running all tests (verbose output)"
	@ echo ""
	@$(TEST_COMMAND) -test.v ./...

build:
	@echo "==> Building executables"
	@ GOOS=linux CGO_ENABLED=0 godep go build -a -installsuffix cgo -ldflags $(LDFLAGS) -o bin/$(EXECUTABLE)-linux $(PKG)
	@ GOOS=darwin CGO_ENABLED=0 godep go build -a -ldflags $(LDFLAGS) -o bin/$(EXECUTABLE)-darwin $(PKG)

install:
	@echo "==> Installing..."
	@godep go install ./...

package: build
	@for p in $(PLATFORMS) ; do \
		echo "==> Tar'ing up $$p/amd64 binary" ; \
		test -f bin/$(EXECUTABLE)-$$p && \
		mv bin/$(EXECUTABLE)-$$p t2 && \
		tar czf $(EXECUTABLE)-$$p.tgz t2 ; \
	done

clean:
	@echo "==> Cleaning up workspace..."
	@go clean ./...
	@rm -rf t2* tnator*.tar.gz
