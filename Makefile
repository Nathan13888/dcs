.DEFAULT_GOAL := install

build:
	# *** BUILDING DCS ***
	go build -o bin/dcs

install:
	make build
	# *** INSTALLING DCS ***
	cp bin/dcs ~/.local/bin/dcs

