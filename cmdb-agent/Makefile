BINARY_NAME=oneagent
VERSION=1.0.0

BUILD_DIR=./build

PLATFORMS=linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64

.PHONY: all clean

all: $(PLATFORMS)

$(PLATFORMS):
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) \
		go build -ldflags "-s -w -X main.Version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-$(word 1,$(subst /, ,$@))-$(word 2,$(subst /, ,$@))$(if $(findstring windows,$(word 1,$(subst /, ,$@))),.exe,) .
	@echo "Built: $(BINARY_NAME)-$(word 1,$(subst /, ,$@))-$(word 2,$(subst /, ,$@))"

clean:
	rm -rf $(BUILD_DIR)
