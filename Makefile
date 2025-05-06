# Makefile is not relevant for golang but I'm lazy to write cmd/shell/main.go everytime I open my ide

.PHONY: hello build run run-bin

hello:
	@echo "Ram Ram mittr"

build:
	go build -o bin/main cmd/shell/main.go

run:
	go run cmd/shell/main.go

run-bin: build
	./bin/main
