#setup > wire > clean > build > run


SERVICE_NAME = order_service_v2
WORKER_MAIN_FILE = server_app
BUILD_DIR = $(PWD)/build
GO= go

setup:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/swaggo/swag/cmd/swag@latest

wire:
	cd internal/ && wire

# clean build file
clean:
	echo "remove bin exe"
	rm -rf $(BUILD_DIR)

# build binary
build:
	echo "build binary execute file"
	make wire
	cd cmd/ && GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(WORKER_MAIN_FILE) .

run:
	make build

	echo "Run service application"
	cd $(BUILD_DIR) && ./$(WORKER_MAIN_FILE)