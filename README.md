# SEAL, COPY & VERIFY - FAST

`godi` stands for "go data integrity" and is a [commandline utility](http://en.wikipedia.org/wiki/Command-line_interface) to generate seal files from a directory tree. This allows to re-check the tree for consistency, and thus verify the data is intact. This is especially useful if valuable, immutable data is retrieved from potentially unreliable media, and copied to another storage device.

As it is very common to verify copy operations, `godi` is able to copy files in the moment is generates the seal, to one or more destinations.

`godi` was inspired by the [media hash list](http://mediahashlist.org) tool, whose seal file format it can read and write. It aims to be easier to use, and faster. On a quad core CPU, `godi` turned to be *6x* the speed of `mhl`.

## Usage Examples

**Protect your valuable, immutable data against silent corruption and create a seal**

```bash
$ godi seal /path/to/data
[...]
Wrote seal file to '/path/to/data/godi_2014-07-29_142701.gobz'
SEAL DONE: WC 7.09s | ->READ ⌗0034 ⌗Δ0004/s ⌰4.70GiB Δ678.61MiB/s | HASH ⌗  9.40GiB Δ  1.33GiB/s
```

**Verify that data is still intact as compared against the seal**

```bash
$ godi verify /path/to/data/godi_2014-07-29_142701.gobz
[...]
VERIFY SUCCESS: None of 35 file(s) changed based on seal in '/path/to/data' [WC  4.794198665s | ->READ ⌗0035 ⌗Δ0007/s ⌰  4.70GiB Δ1004.33MiB/s | HASH ⌗  9.40GiB Δ  1.96GiB/s]
```

**Seal valuable data and duplicate it to two backup devices**

```bash
$ godi sealed-copy /path/to/data -- /media/d1 /media/d2
[...]
Wrote seal file to '/media/d1/godi_2014-07-29_144842.gobz'
Wrote seal file to '/media/d2/godi_2014-07-29_144842.gobz'
SEAL DONE: WC  3.761588954s |   ->READ ⌗0475 ⌗Δ0126/s ⌰407.90MiB Δ108.44MiB/s |   HASH ⌗815.81MiB Δ216.88MiB/s |   WRITE ⌗0950 ⌗Δ0252/s ⌰815.81MiB Δ216.88MiB/s (16 skipped)
```

**Seal valuable data and duplicate it to two backup devices [PARANOID VERSION]**

```bash
$ godi sealed-copy --verify /path/to/data -- /media/d3 /media/d4
[...]
Wrote seal file to '/media/d3/godi_2014-07-29_145114.gobz'
Wrote seal file to '/media/d4/godi_2014-07-29_145114.gobz'
SEAL DONE: WC  3.039262074s |   ->READ ⌗0475 ⌗Δ0156/s ⌰407.90MiB Δ134.21MiB/s |   HASH ⌗815.81MiB Δ268.42MiB/s |   WRITE ⌗0950 ⌗Δ0312/s ⌰815.81MiB Δ268.42MiB/s (16 skipped)
[...]
VERIFY SUCCESS: None of 475 file(s) changed based on seal in '/media/d3'
VERIFY SUCCESS: None of 475 file(s) changed based on seal in '/media/d4' [WC  2.297587367s |   ->READ ⌗0950 ⌗Δ0413/s ⌰815.81MiB Δ355.07MiB/s |   HASH ⌗  1.59GiB Δ710.15MiB/s]
```

Read more [on the documentation page](http://byron.github.io/godi)

## Installation

* Download and extract the multi-platform archive from the [latest release](https://github.com/Byron/godi/releases)
* Typing `godi/<platform>/godi` in a [terminal](http://en.wikipedia.org/wiki/Terminal_emulator) from the extraction point displays help on the respective platform.

If you are using godi more regularly, it is adviced to put it into your shell's [executable search path](http://en.wikipedia.org/wiki/PATH_(variable))

## Features

* protect your data against unnoticed and silent corruption
* [atomically](http://en.wikipedia.org/wiki/Atomic_operation) duplicate data to one or more locations while protecting it against corruption
* use all available device bandwidth and CPU cores to finish your work *much* faster.
* runs on all major and minor platforms

## Development Status

[![Build Status](https://travis-ci.org/Byron/godi.svg?branch=master)](https://travis-ci.org/Byron/godi)

Latest releases [can be found here](https://github.com/Byron/godi/releases).

## Credits

`godi` uses the following libraries:

* [codegangsta/cli](https://github.com/codegangsta/cli)

`godi` is inspired by

* [media hash list](http://mediahashlist.org)

### LICENSE

This open source software is licensed under [GNU Lesser General Public License](https://github.com/Byron/godi/blob/master/LICENSE.md)
