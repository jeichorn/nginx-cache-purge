#!/bin/bash
set -e
export GOPATH=`realpath .`
go get github.com/go-fsnotify/fsnotify
go get github.com/koyachi/go-term-ansicolor/ansicolor
go get gopkg.in/redis.v3
go get github.com/jeichorn/go-caching
cd $GOPATH
cd src/nginxcp
go build -v
cd $GOPATH
cd src/nginx-cache-purge
go install
