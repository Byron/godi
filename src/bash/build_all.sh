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
# Make sure we keep possible debug versions around ... 
mv $bindatadest $bindatadest.orig &>/dev/null
go-bindata -ignore "[/\\\\]\\..*" -nomemcopy=true -o $bindatadest -pkg=server -prefix=web/client/app web/client/app/... || exit  $?
go fmt $bindatadest || exit $?

gox -ldflags="-w -s" -tags web -verbose -arch=amd64 -os="linux darwin windows" -output="build/{{.OS}}_{{.Arch}}/{{.Dir}}" github.com/Byron/godi

# Restore possibly existing debug version
mv $bindatadest.orig $bindatadest &>/dev/null
