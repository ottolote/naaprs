VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -X main.Version=$(VERSION)

GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

CMD_DIR = cmd
BIN_DIR = bin
BINARY_NAME = naaprs
DOCKER_REGISTRY = ottokl
DOCKER_REPO = $(DOCKER_REGISTRY)/$(BINARY_NAME)
DOCKER_IMAGE = $(DOCKER_REPO):$(VERSION)

PLATFORMS = linux/amd64 linux/386 linux/arm/v6 linux/arm/v7 linux/arm64
BUILD_TARGETS = $(addprefix build-, $(PLATFORMS))
PUSH_TARGETS = $(addprefix push-, $(PLATFORMS))

# Main targets
build: clean $(BIN_DIR)/$(BINARY_NAME)

$(BIN_DIR)/$(BINARY_NAME):
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $@ $(CMD_DIR)/naaprs/main.go

push: $(BUILD_TARGETS) $(PUSH_TARGETS) docker-push-multi-arch

push-latest: push docker-push-multi-arch-latest

# Docker build and push targets
$(BUILD_TARGETS):
	platform=$$(echo $@ | sed 's/build-//'); \
	docker buildx build \
	--build-arg LDFLAGS="$(LDFLAGS)" \
	--build-arg GOOS=$$(echo $$platform | cut -d/ -f1) \
	--build-arg GOARCH=$$(echo $$platform | cut -d/ -f2) \
	--build-arg GOARM=$$(echo $$platform | cut -d/ -f3 | sed 's/v//') \
	--platform $$platform -t $(DOCKER_IMAGE)-$$(echo $$platform | tr '/' '-') --load .

$(PUSH_TARGETS): push-%: build-%
	docker push $(DOCKER_IMAGE)-$$(echo $* | tr '/' '-')

docker-push-multi-arch: $(PUSH_TARGETS)
	docker manifest create $(DOCKER_IMAGE) $(foreach platform,$(PLATFORMS),--amend $(DOCKER_IMAGE)-$$(echo $(platform) | tr '/' '-'))
	docker manifest push $(DOCKER_IMAGE)

docker-push-multi-arch-latest: $(PUSH_TARGETS)
	docker manifest create $(DOCKER_REPO):latest $(foreach platform,$(PLATFORMS),--amend $(DOCKER_IMAGE)-$$(echo $(platform) | tr '/' '-'))
	docker manifest push $(DOCKER_REPO):latest

clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/*

test:
	$(GOTEST) ./...

deps:
	$(GOGET) ./...

run:
	./$(BIN_DIR)/$(BINARY_NAME)

.PHONY: build clean test deps run push push-latest docker-push-multi-arch docker-push-multi-arch-latest $(BUILD_TARGETS) $(PUSH_TARGETS)

