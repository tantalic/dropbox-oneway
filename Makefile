PKG := "tantalic.com/dropbox-oneway"
GOVERSION := 1.8.0
DOCKER_IMAGE := "tantalic/dropbox-oneway"

COMMIT = $(strip $(shell git rev-parse --short HEAD))
VERSION := $(strip $(shell git describe --always --dirty))


.PHONY: run build linux-amd64 linux macos freebsd openbsd netbsd dragonfly windows docker-build docker-push binaries update-ca help
.DEFAULT_GOAL := help	

run: ## Run from source, using the local go installation
	go run *.go

build:  ## Create a build for the current platform using local go installation
	go build -o build/dropbox-oneway

linux-amd64:
	docker run --env GOOS=linux --env GOARCH=amd64 --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/dropbox-oneway-linux_amd64

linux: linux-amd64
	docker run --env GOOS=linux --env GOARCH=386                 --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/dropbox-oneway-linux_386
	docker run --env GOOS=linux --env GOARCH=arm   --env GOARM=6 --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/dropbox-oneway-linux_armv6
	docker run --env GOOS=linux --env GOARCH=arm   --env GOARM=7 --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/dropbox-oneway-linux_armv7
	docker run --env GOOS=linux --env GOARCH=arm64               --env CGO_ENABLED=0 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -a -tags netgo -ldflags '-w' -o build/dropbox-oneway-linux_arm64

macos:
	docker run --env GOOS=darwin --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-darwin_amd64	

freebsd:
	docker run --env GOOS=freebsd --env GOARCH=amd64 --env         --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-freebsd_amd64
	docker run --env GOOS=freebsd --env GOARCH=386   --env         --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-freebsd_386
	docker run --env GOOS=freebsd --env GOARCH=arm   --env GOARM=6 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-freebsd_armv6
	docker run --env GOOS=freebsd --env GOARCH=arm   --env GOARM=7 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-freebsd_armv7

openbsd:
	docker run --env GOOS=openbsd --env GOARCH=amd64 --env         --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-openbsd_amd64
	docker run --env GOOS=openbsd --env GOARCH=386   --env         --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-openbsd_386
	docker run --env GOOS=openbsd --env GOARCH=arm   --env GOARM=6 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-openbsd_armv6
	docker run --env GOOS=openbsd --env GOARCH=arm   --env GOARM=7 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-openbsd_armv7

netbsd:
	docker run --env GOOS=netbsd --env GOARCH=amd64 --env         --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-netbsd_amd64
	docker run --env GOOS=netbsd --env GOARCH=386   --env         --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-netbsd_386
	docker run --env GOOS=netbsd --env GOARCH=arm   --env GOARM=6 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-netbsd_armv6
	docker run --env GOOS=netbsd --env GOARCH=arm   --env GOARM=7 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-netbsd_armv7

dragonfly:
	docker run --env GOOS=dragonfly --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-dragonfly_amd64

windows:
	docker run --env GOOS=windows --env GOARCH=amd64 --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-windows_amd64
	docker run --env GOOS=windows --env GOARCH=386   --rm -v "${PWD}":/go/src/$(PKG) -w /go/src/$(PKG) golang:$(GOVERSION) go build -o build/dropbox-oneway-windows_386

docker-image: linux-amd64 ## Build a docker image
	docker build \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg VERSION=$(VERSION) \
		--build-arg VCS_REF=$(COMMIT) \
		-t $(DOCKER_IMAGE):$(VERSION) .

docker-push: ## Push the docker image to DockerHub
	docker push $(DOCKER_IMAGE):$(VERSION)

binaries: linux macos freebsd openbsd netbsd dragonfly windows ## Create binaries for all supported platforms

update-ca: ## Download the latest CA roots
	curl --time-cond certs/ca-certificates.crt -o certs/ca-certificates.crt https://curl.haxx.se/ca/cacert.pem

help: ## Print available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
