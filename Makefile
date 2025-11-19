dirs = ./interfaces/... ./commands/... ./services/... ./bot/... .

VERSION != git describe --tags | sed 's/\([^-]*-\)g/r\1/'

# packr2 has incompatibility issue if used without GO111MODULE=on, see:
# https://github.com/gobuffalo/packr/issues/183
# https://github.com/gobuffalo/packr/pull/265
export GO111MODULE=on

.PHONY: all all-static
all: build ## Default action. Compile resources and builds MumbleDJ.
all-static: build-static

.PHONY: build
build: *.go  ## Builds MumbleDJ.
	@env go build -ldflags '-X "main.version=$(VERSION)"' .

.PHONY: build-static
build-static: *.go
	@env go build -a -buildmode=pie -tags netgo,osusergo,static_build -trimpath -ldflags '-s -w -X "main.version=$(VERSION)" -extldflags "-fno-pic -static"' -o mumbledj .

.PHONY: test
test: ## Runs unit tests for MumbleDJ.
	@env go test $(dirs)

.PHONY: coverage
coverage: ## Runs coverage tests for MumbleDJ.
	@env overalls -project=github.com/leoverto/mumbledj -covermode=atomic
	@mv overalls.coverprofile coverage.txt

.PHONY: clean
clean: ## Removes compiled MumbleDJ binaries.
	@rm -f mumbledj*

.PHONY: install
install: ## Copies MumbleDJ binary to /usr/local/bin for easy execution.
	@cp -f mumbledj* /usr/local/bin/mumbledj

# .PHONY: dist
# dist: ## Performs cross-platform builds via gox for multiple Linux platforms.
# 	@go install -u github.com/mitchellh/gox
# 	@gox -cgo -osarch="linux/amd64 linux/386"

.PHONY: help
help: ## Shows this helptext.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
