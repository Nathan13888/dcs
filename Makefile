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

publish:
	make publish-ghcr

publish-ghcr:
	make docker-build
	# TODO: specify tag version
	docker tag dcs:latest docker.pkg.github.com/nathan13888/dcs/dcs:latest
	docker push docker.pkg.github.com/nathan13888/dcs/dcs:latest

pull-ghcr:
	docker pull docker.pkg.github.com/nathan13888/dcs/dcs:latest

install:
	# *** INSTALLING DCS ***
	cp bin/dcs ~/.local/bin/dcs

test:
	go test -v ./...

