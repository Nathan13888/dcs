#!/bin/bash

make build
pkill dcs
make install
