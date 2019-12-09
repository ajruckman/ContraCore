#!/bin/bash

cd $GOPATH/src/github.com/ajruckman/ContraCore/internal/provision || exit
go generate

rm -rf $GOPATH/pkg/mod/github.com/ajruckman/!contra!core@v0.0.1/*
mkdir -p $GOPATH/pkg/mod/github.com/ajruckman/!contra!core@v0.0.1/
cp -r  $GOPATH/src/github.com/ajruckman/ContraCore/* $GOPATH/pkg/mod/github.com/ajruckman/!contra!core@v0.0.1/

pkill -9 coredns

cd $GOPATH/src/github.com/coredns/coredns/ || exit

go env -w GOPRIVATE="github.com/ajruckman/ContraCore,github.com/ajruckman/xlib"
sed -i '/ContraCore/d' go.sum
make || exit

./coredns -conf $GOPATH/src/github.com/ajruckman/ContraCore/Corefile
