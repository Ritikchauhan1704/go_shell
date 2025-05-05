# Makefile is not relevant for golang but I'm lazy to write app/main.go everytime I open my ide

hello:
	echo "Ram Ram mittr"

build:
	go build -o bin/main app/main.go

run:
	go run app/main.go