.DEFAULT_GOAL := install

build:
	# *** BUILDING DCS ***
	go build -o bin/dcs -ldflags "-X 'dcs/config.BuildUser=$$(id -u -n)' -X 'dcs/config.BuildTime=$$(date)' -s -w"

install:
	make build
	# *** INSTALLING DCS ***
	cp bin/dcs ~/.local/bin/dcs

test:
	go test -v ./...

