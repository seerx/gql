#!/bin/bash
export GOOS=linux
export CGO_ENABLE=0

# 编译健康监测程序
go build -o watcher-linux-amd64; echo build `pwd`; cd ..

# 编译主程序
# cd test; go get; go build -o test-linux-amd64; echo build `pwd`; cd ..

export GOOS=darwin