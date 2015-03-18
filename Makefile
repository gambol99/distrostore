#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-03-16 14:28:12 +0000 (Mon, 16 Mar 2015)
#
#  vim:ts=2:sw=2:et
#
NAME=cluster-store
AUTHOR=gambol99
VERSION=$(shell awk '/const Version/ { print $$4 }' version.go | sed 's/"//g')

default: build

build:
	go get
	go build

test: build
	go test -v

.PHONY: build release changelog
