.PHONY: setup run nuke help

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  setup     : Downloads the Go module dependencies"
	@echo "  run       : Builds and runs the project using docker-compose"
	@echo "  nuke      : Stops and removes all project containers, networks, and volumes"
	@echo "  help      : Displays this help message"

setup:
	go mod download

run:
	docker-compose up --build

nuke:
	docker-compose down -v --remove-orphans
	docker volume prune -f
