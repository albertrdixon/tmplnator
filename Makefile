PROJECT = github.com/albertrdixon/tmplnator
TEST_COMMAND = go test
EXECUTABLE = t2
PKG = cmd/t2/t2.go
LDFLAGS = -s
PLATFORMS = linux darwin
BUILD_ARGS = ""
TOOLS = glide

.PHONY: save restore test test-verbose build install package clean

all: test install

help:
	@echo "Available targets:"
	@echo ""
	@echo "  save"
	@echo "  restore"
	@echo "  test"
	@echo "  test-verbose"
	@echo "  build"
	@echo "  install"
	@echo "  package"
	@echo "  clean"

tools:
	go get -u -v -ldflags -s github.com/Masterminds/glide

save:
	@echo "---> Saving dependencies..."
	@glide update

restore:
	@echo "---> Restoring dependencies..."
	@glide install

test:
	@echo "---> Running all tests"
	@echo ""
	@$(TEST_COMMAND) .

test-verbose:
	@echo "---> Running all tests (verbose output)"
	@ echo ""
	@$(TEST_COMMAND) -test.v .

build:
	@echo "---> Building executables"
	@ GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags $(LDFLAGS) -o bin/$(EXECUTABLE)-linux $(PKG)
	@ GOOS=darwin CGO_ENABLED=0 go build -a -ldflags $(LDFLAGS) -o bin/$(EXECUTABLE)-darwin $(PKG)

install:
	@echo "---> Installing..."
	@CGO_ENABLED=0 go install -a -ldflags $(LDFLAGS) $(PKG)

package: build
	@for p in $(PLATFORMS) ; do \
		echo "---> Tar'ing up $$p/amd64 binary" ; \
		test -f bin/$(EXECUTABLE)-$$p && \
		cp bin/$(EXECUTABLE)-$$p t2 && \
		tar czf $(EXECUTABLE)-$$p.tgz t2 ; \
	done

clean:
	@echo "---> Cleaning up workspace..."
	@go clean ./...
	@rm -rf t2* tnator*.tar.gz
