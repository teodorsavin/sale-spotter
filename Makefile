.PHONY: setup run nuke help

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  setup     : Downloads the Go module dependencies"
	@echo "  run       : Builds and runs the project using docker-compose"
	@echo "  nuke      : Stops and removes all project containers, networks, and volumes"
	@echo "  help      : Displays this help message"

setup:
	@if ! go version | grep -q "go1\.[2-9]\|go[2-9][0-9]"; then \
		echo "Go version 1.19 or higher is required. Please install it. The easiest way is to use asdf"; \
		exit 1; \
	fi
	go mod download

run:
	docker-compose up --build

nuke:
	docker-compose down -v --remove-orphans
	docker volume prune -f
