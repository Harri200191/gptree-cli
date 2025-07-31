APP_NAME = gptree
BUILD_DIR = .build
SRC = main.go

# Detect OS for appropriate binary extension
ifeq ($(OS),Windows_NT)
    EXT = .exe
    RM = powershell -Command "Remove-Item -Recurse -Force $(BUILD_DIR)"
else
    EXT =
    RM = rm -rf $(BUILD_DIR)
endif

all: build

build:
	@echo "🔧 Building $(APP_NAME)..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME)$(EXT) $(SRC)
	@echo "✅ Built $(BUILD_DIR)/$(APP_NAME)$(EXT)"

clean:
	@echo "🧹 Cleaning..."
	-@$(RM)
	@echo "✅ Cleaned up build directory."

run:
	@$(BUILD_DIR)/$(APP_NAME)$(EXT) . --help

install:
	sudo cp $(BUILD_DIR)/$(APP_NAME)$(EXT) /usr/local/bin/$(APP_NAME)
	sudo chmod +x /usr/local/bin/$(APP_NAME)
	@echo "🚀 Installed $(APP_NAME) to /usr/local/bin"
