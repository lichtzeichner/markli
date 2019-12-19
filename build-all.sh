#!/usr/bin/env -S bash -e

function build() {
    local os=$1
    local arch=$2
    local suffix=$3
    filename="markli-${os}-${arch}${suffix}"
    echo "Building ${filename}"
    GOOS=$os GOARCH=$arch go build -o dist/${filename} 
}

build windows 386 exe
build windows amd64 exe
build linux 386
build linux amd64
build linux arm
build linux amd64
build darwin amd64