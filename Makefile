.DEFAULT_GOAL := build-n-install

update:
	./update.sh

build-n-install:
	make build
	make install

build:
	# *** BUILDING DCS ***
	go build -o bin/dcs -ldflags "-X 'dcs/config.BuildUser=$$(id -u -n)' -X 'dcs/config.BuildTime=$$(date)' -s -w"

docker-build:
	docker build -t dcs .

install:
	# *** INSTALLING DCS ***
	cp bin/dcs ~/.local/bin/dcs

test:
	go test -v ./...

