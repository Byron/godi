#!/bin/bash

if ![[ `which go-bindata` &>/dev/null ]]; then
	echo "go-bindata couldn't be foudn in PATH, go get github.com/jteeuwen/go-bindata/... should do the job"
	exit 2
fi

version=${1:?First argument must be the version you are building, like v1.1.0}

base="`dirname $0`/../.."
build_dir=build

cd $base || exit $?

bindatadest=web/server/clientdata.go
# NOTE: use -debug flag to compile in debug mode instead, to keep files live !
# go-bindata -debug ...
# Make sure we keep possible debug versions around ... 
mv $bindatadest $bindatadest.orig &>/dev/null
go-bindata -ignore "[/\\\\]\\..*" -nomemcopy=true -o $bindatadest -pkg=server -prefix=web/client/app web/client/app/... || exit  $?
go fmt $bindatadest || exit $?

for mode in web ""; do
	arg=""
	suffix=""
	if [[ $mode == web ]]; then
		arg="-tags web"
		suffix="_"
	fi
	echo "building $mode godi ..."
	basename=godi_${mode}${suffix}${version}
	gox -ldflags="-w -s" $arg -verbose -arch=amd64 -os="linux darwin windows" -output="$build_dir/$basename/{{.OS}}_{{.Arch}}/{{.Dir}}" github.com/Byron/godi

	(cd $build_dir && zip --quiet -r -9 $basename.zip $basename)&
done

echo "Waiting for archives to be created ..."
wait

# Restore possibly existing debug version
mv $bindatadest.orig $bindatadest &>/dev/null
