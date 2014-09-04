#!/bin/bash

if ![[ `which go-bindata` &>/dev/null ]]; then
	echo "go-bindata couldn't be foudn in PATH, go get github.com/jteeuwen/go-bindata/... should do the job"
	exit 2
fi

base="`dirname $0`/../.."

cd $base || exit $?

bindatadest=web/server/clientdata.go
# NOTE: use -debug flag to compile in debug mode instead, to keep files live !
# go-bindata -debug ...
go-bindata -ignore "[/\\\\]\\..*" -nomemcopy=true -o $bindatadest -pkg=server -prefix=web/client/app web/client/app/... || exit  $?
go fmt $bindatadest || exit $?

gox -ldflags="-w -s" -verbose -arch=amd64 -os="linux darwin windows" -output="build/{{.OS}}_{{.Arch}}/{{.Dir}}" github.com/Byron/godi

# Just make sure we keep a debug-compatible version, to safe some space in our repository
git checkout $bindatadest
