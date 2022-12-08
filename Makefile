.PHONY: all main clean test cover lint

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
MAIN_DIR					:= ./...
BUILD_DIR   			:= ./build
COVERAGE_PROFILE	:= cover.prof
BIN_FILE					:= $(BUILD_DIR)/sample
LINT_CONFIG				:= review_config.yml

### PHONY ターゲットのビルドルール

all: clean test main

clean:
	rm -fr $(BUILD_DIR) $(COVERAGE_PROFILE)
main:
	$(GO_BUILD) -o $(BIN_FILE) $(MAIN_DIR)
test:
	$(GO_TEST) $(MAIN_DIR)
cover:
	$(GO_TEST) $(MAIN_DIR) -covermode=count -coverprofile=$(COVERAGE_PROFILE)
	$(GO_COVER) -func=$(COVERAGE_PROFILE)
lint:
	$(GOLINT) -c $(LINT_CONFIG)
