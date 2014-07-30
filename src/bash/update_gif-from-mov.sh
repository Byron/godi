#!/usr/bin/env bash

if ![ `which ffmpeg` &>/dev/null ] || ![ `which convert` &>/dev/null ]; then
	echo "Couldn't find 'ffmpeg' or 'convert' program in PATH"
	exit 2
fi

base=`dirname $0`

for smov in `find $base/../mov -name "*.mov" -type f`; do
	dest=$base/../../lib/gif/`basename $smov`.gif
	mkdir -p `dirname $dest`
	echo "FFMPEG: $smov -> $dest"
	ffmpeg  -loglevel error -y -i $smov -pix_fmt rgb24 -r 5 $dest || exit $?
	echo "OPTIMIZE: $dest"
	convert -layers Optimize $dest $dest
done || exit $?
