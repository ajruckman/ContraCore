#!/bin/bash

set +H
git config --global --add url."git@github.com:".insteadOf "https://github.com/"

go env -w GOARCH="amd64"
go env -w GOOS="linux"

cd $GOPATH/src/github.com/ajruckman/ContraCore/internal || exit
go generate

rm -rf $GOPATH/pkg/mod/github.com/ajruckman/!contra!core@v0.0.1/*
mkdir -p $GOPATH/pkg/mod/github.com/ajruckman/!contra!core@v0.0.1/
cp -r  $GOPATH/src/github.com/ajruckman/ContraCore/* $GOPATH/pkg/mod/github.com/ajruckman/!contra!core@v0.0.1/

pkill -9 coredns

cd $GOPATH/src/github.com/coredns/coredns/ || exit

go env -w GOPRIVATE="github.com/ajruckman/ContraCore,github.com/ajruckman/xlib"

#go env -w GOARCH="arm"
#go env -w GOARM="7"

sed -i '/ContraCore/d' go.sum
make || exit

./coredns -conf $GOPATH/src/github.com/ajruckman/ContraCore/Corefile

