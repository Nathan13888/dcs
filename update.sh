#!/bin/sh

git pull
make build
pkill dcs
make install
dcs version
