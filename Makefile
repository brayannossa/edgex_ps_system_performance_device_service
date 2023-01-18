.PHONY: build test clean docker

GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-os/device-os
.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
DOCKER_TAG=$(VERSION)-dev

GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-os.Version=$(VERSION)"
GOTESTFLAGS?=-race

GIT_SHA=$(shell git rev-parse HEAD)

build: $(MICROSERVICES)
	$(GOCGO) install -tags=safe

cmd/device-os/device-os:
	go mod tidy
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/device-os

docker:
	docker build \
		-f example/cmd/device-os/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/device-os:$(GIT_SHA) \
		-t edgexfoundry/device-os:$(DOCKER_TAG) \
		.

test:
	go mod tidy
	GO111MODULE=on go test $(GOTESTFLAGS) -coverprofile=coverage.out ./...
	GO111MODULE=on go vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	./bin/test-attribution-txt.sh
	./bin/test-go-mod-tidy.sh

clean:
	rm -f $(MICROSERVICES)
