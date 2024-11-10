# Variables
CMD_DIR := cmd/server
APP_NAME := trinity
SCRIPT_DIR := scripts

# Default target
all: build

# Run the application
run:
	cd $(CMD_DIR) && go run .

# Build the application
build:
	cd $(CMD_DIR) && go build -o $(APP_NAME)

# Clean the generated binaries
clean:
	rm -f $(CMD_DIR)/$(APP_NAME)

# Swagger
swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g $(CMD_DIR)/main.go -o ./docs

# Run all steps from the script
run-all:
	bash $(SCRIPT_DIR)/run-all.sh

# Help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo
	@echo "Usage:"
	@echo "  make run         Run the application"
	@echo "  make swag		  Run the swagger"
	@echo "  make build       Build the application"
	@echo "  make clean       Clean the generated binaries"
	@echo "  make help        Show this help message"
