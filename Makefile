.PHONY: all build build-% clean test test-% cover lint

### コマンドの定義
GO          			:= go
GO_BUILD    			:= $(GO) build
GO_FORMAT   			:= $(GO) fmt
GOFMT       			:= gofmt
GO_LIST     			:= $(GO) list
GOLINT      			:= golangci-lint run -v
GO_TEST     			:= $(GO) test -v
GO_VET      			:= $(GO) vet
GO_LDFLAGS  			:= -ldflags="-s -w"
GO_TOOL						:= $(GO) tool
GO_COVER					:= $(GO_TOOL) cover 

### ターゲットパラメータ
NAMES							:= logging registry
MODULE_NAME				:= github.com/NARH/go.tools
MAIN_DIR					:= $(MODULE_NAME)/...
BUILD_DIR   			:= ./build
COVERAGE_PROFILE	:= cover.prof
BIN_FILE					:= $(BUILD_DIR)/sample
LINT_CONFIG				:= review_config.yml

### PHONY ターゲットのビルドルール

all: clean test build

clean:
	rm -fr $(BUILD_DIR) $(COVERAGE_PROFILE)

build: $(addprefix build-, $(NAMES))

build-%:
	$(GO_BUILD) $(MODULE_NAME)/${@:build-%=%}

test: $(addprefix test-, $(NAMES))

test-%:
	$(GO_TEST) $(MODULE_NAME)/${@:test-%=%}

cover: $(addprefix cover-, $(NAMES))

cover-%:
	$(GO_TEST) $(MODULE_NAME)/${@:cover-%=%} -covermode=count -coverprofile=${@:cover-%=%}-$(COVERAGE_PROFILE)
	$(GO_COVER) -func=${@:cover-%=%}-$(COVERAGE_PROFILE)

lint:
	$(GOLINT) -c $(LINT_CONFIG)

