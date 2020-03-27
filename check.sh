#/bin/bash

set -e

$GOPATH/bin/golint -set_exit_status ./src/...
$GOPATH/bin/gosec ./...
$GOPATH/bin/gocheckstyle -config=.go_style src
$GOPATH/bin/go-namecheck -rules .go_naming_rules ./...
$GOPATH/bin/dupl -t 200

go get -d -v ./...
go test -cover ./...
