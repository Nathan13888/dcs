#!/bin/bash

git pull
make build
pkill dcs
make install
