##=======================================================================##
## Makefile
## Created: Wed Aug 05 14:35:14 PDT 2015 @941 /Internet Time/
# :mode=makefile:tabSize=3:indentSize=3:
## Purpose:
##======================================================================##

SHELL=/bin/bash
PROJECT_NAME = tileserver
GPATH = $(shell pwd)

.PHONY: fmt install get-deps scrape build clean

install: fmt get-deps
	./install.sh
	@GOPATH=${GPATH} go build -o tile_server ${PROJECT_NAME}/main.go

fmt:
	@GOPATH=${GPATH} gofmt -s -w ${PROJECT_NAME}

get-deps:
	@GOPATH=${GPATH} go get -v github.com/cihub/seelog
	@GOPATH=${GPATH} go get -v github.com/gorilla/mux

scrape:
	@find src -type d -name '.hg' -or -type d -name '.git' | xargs rm -rf

clean:
	@GOPATH=${GPATH} go clean
