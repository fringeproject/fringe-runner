NAME := fringe-runner
PACKAGE_NAME := fringe-runner
DOCKER_REPO := fringeproject
VERSION := beta
REVISION := $(shell git rev-parse --short=8 HEAD || echo unknown)
BRANCH := $(shell git show-ref | grep "$(REVISION)" | grep -v HEAD | awk '{print $$2}' | sed 's|refs/remotes/origin/||' | sed 's|refs/heads/||' | sort | head -n 1)
BUILT := $(shell date -u +%Y-%m-%dT%H:%M:%S%z)

PKG = github.com/fringeproject/$(PACKAGE_NAME)
COMMON_PACKAGE_NAMESPACE=$(PKG)/common

# Update common/version.go constants to matchs the current repository status
GO_LDFLAGS ?= -X $(COMMON_PACKAGE_NAMESPACE).NAME=$(PACKAGE_NAME) \
			  -X $(COMMON_PACKAGE_NAMESPACE).VERSION=$(VERSION) \
			  -X $(COMMON_PACKAGE_NAMESPACE).REVISION=$(REVISION) \
			  -X $(COMMON_PACKAGE_NAMESPACE).BUILT=$(BUILT) \
			  -X $(COMMON_PACKAGE_NAMESPACE).BRANCH=$(BRANCH)

GO_BUILD_OS ?= darwin linux windows
GO_BUILD_ARCH ?= 386 amd64
BUILD_FOLDER := ./build


update:
	@go get -u; \
	go mod tidy -v; \
	echo "Update done."

clean:
	@go clean ./... ; \
	echo "Clean done."

lint:
	@go get -u github.com/golangci/golangci-lint@master ; \
	golangci-lint run ./... ; \
	go mod tidy ; \
	echo "Lint Done."

fmt:
	@go fmt ./...; \
	echo "Fmt done."

GOX:
	@go get github.com/mitchellh/gox; \
	echo "Install gox done."

build: GOX
	@mkdir -p $(BUILD_FOLDER); \
	gox -os="$(GO_BUILD_OS)" -arch="$(GO_BUILD_ARCH)" -ldflags "$(GO_LDFLAGS)" -output="$(BUILD_FOLDER)/$(NAME)-{{.OS}}-{{.Arch}}" $(PKG); \
	echo "Build done."

install:
	@go install; \
	echo "Install done."

upload: build
	@scp ./resources/*.json fringe:/opt/fringe-static/fringe-runner/resources/; \
	scp $(BUILD_FOLDER)/* fringe:/opt/fringe-static/fringe-runner/builds/; \
	echo "Upload done."

docker-build:
	@docker build -t "$(NAME)" . ; \
	echo "Docker build done."

all: clean fmt update lint
