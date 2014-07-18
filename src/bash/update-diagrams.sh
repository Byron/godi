#!/usr/bin/env bash

if ![ `which ditaa` &>/dev/null ]; then
	echo "Couldn't find 'ditaa' program in PATH"
	exit 2
fi

base=`dirname $0`

for fdia in `find $base/../md -name "*.dia.txt" -type f`; do
	dest=$base/../md/img/`basename $fdia`.png
	mkdir -p `dirname $dest`
	ditaa $fdia $dest || exit $?
done || exit $?