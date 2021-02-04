.DEFAULT_GOAL := run

run:
	go run main.go

build:
	go build -o bin/dcs

install:
	cp bin/dcs $HOME/.local/bin/dcs

