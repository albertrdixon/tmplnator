PROJECT = github.com/albertrdixon/tmplnator
TEST_COMMAND = godep go test
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

build:
	@echo "==> Building executables"
	@gox -osarch="linux/amd64 darwin/amd64" -output="{{.Dir}}-{{.OS}}-{{.Arch}}" ./...

install:
	@echo "==> Installing..."
	@godep go install ./...

package: build
	@for p in $(PLATFORMS) ; do \
		echo "==> Tar'ing up $$p/amd64 binary" ; \
		test -f t2-$$p-amd64 && mv t2-$$p-amd64 t2 && tar czf tnator-$$p-amd64.tar.gz t2 ; \
		rm -f t2 ; \
	done

clean:
	go clean ./...
	rm -rf t2* tnator*.tar.gz
