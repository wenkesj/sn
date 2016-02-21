.PHONY: all deps install test

all: deps install test

build: deps install

deps:
	go get -u github.com/gonum/plot
	go get -u github.com/garyburd/redigo/redis

install:
	go install github.com/wenkesj/sn/plots
	go install github.com/wenkesj/sn/sn
	go install github.com/wenkesj/sn/net
	go install github.com/wenkesj/sn/group
	go install github.com/wenkesj/sn/vars
	go install github.com/wenkesj/sn/sim

test:
	go test -cover -v ./tests/
