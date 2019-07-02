NAME=botzito
VERSION=0.1-rc2
OS ?= linux
PKG ?= github.com/guilhermebr/botzito/cmd

.PHONY: compile
compile:
	@echo "==> Building the project"
	@env GOOS=${OS} GOARCH=amd64 go build -v -o cmd/${NAME} ${PKG}

.PHONY: build
build: clean compile
	@echo "==> Building the docker image"
	@docker build -t guilhermebr/botzito:latest cmd -f cmd/Dockerfile

.PHONY: push
push: build
	@echo "==> Pushing to registry"
	@docker push guilhermebr/botzito:latest

.PHONY: clean
clean:
	@echo "==> Cleaning releases"
	@GOOS=${OS} go clean -i -x ./...
