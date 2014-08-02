# SEAL, COPY & VERIFY - FAST

`godi` stands for "go data integrity" and is a [commandline utility](http://en.wikipedia.org/wiki/Command-line_interface) to generate seal files from a directory tree. This allows to re-check the tree for consistency, and thus verify the data is intact. This is especially useful if valuable, immutable data is retrieved from potentially unreliable media, and copied to another storage device.

As it is very common to verify copy operations, `godi` is able to copy files in the moment is generates the seal, to one or more destinations.

`godi` was inspired by the [media hash list](http://mediahashlist.org) tool, whose seal file format it can read and write. It aims to be easier to use, and faster. On a quad core CPU, `godi` turned to be *6x* the speed of `mhl`.

## Why Godi

`godi` helps you to protect your data against unnoticed corruption, which represents an undesired change. Data files suitable for such protection is

* **valuable** and reproducing it is either costly or prohibitive
    + Imagine you are on set of a block buster, and have to handle the camera data of a multi-million action shot that took weeks to prepare.
    + Your digital product is delivered to the client, and you want him to be sure that data was not corrupted in-flight. It's also an insurance for you as you can prove that the data was still intact on your end.
* **immutable** and will not be changed directly
    + You took pictures of the Niagara Falls, and even though you want to do some post-processing, you will always keep the original. It should remain exactly as is, and will only serve as source for various image manipulations.

Therefore, if your data files change a lot, `godi` is not for you.

## Usage Examples

**Protect your valuable, immutable data against silent corruption and create a seal**
![seal](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_seal.mov.gif)

**Verify that data is still intact as compared against the seal**
![verify](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_verify.mov.gif)


**Seal valuable data and duplicate it to two backup devices**
![sealed-copy](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_sealed-copy.mov.gif)

**Seal valuable data and duplicate it to two backup devices [PARANOID VERSION]**
![sealed-copy-verify](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_sealed-copy-verify.mov.gif)

Read more [on the documentation page](http://byron.github.io/godi)

## Features

* protect your data against unnoticed and silent corruption
* [atomically](http://en.wikipedia.org/wiki/Atomic_operation) duplicate data to one or more locations while protecting it against corruption, without fearing to overwrite existing files.
* use all [available device bandwidth and CPU cores](details.md#Performance Considerations) to finish your work *much* faster
* runs on all major and minor [platforms and architectures](http://golang.org/doc/install#requirements)

## Installation

* Download and extract the multi-platform archive from the [latest release](https://github.com/Byron/godi/releases)
* Typing `godi/<platform>/godi` in a [terminal](http://en.wikipedia.org/wiki/Terminal_emulator) from the extraction point displays help on the respective platform.

If you are using godi more regularly, it is adviced to put it into your shell's [executable search path](http://en.wikipedia.org/wiki/PATH_(variable))

## Videos

[![Godi Performance Video](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/png/godi_performance-thumb.png)](https://vimeo.com/102326726)
[![Godi Usage Video](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/png/godi_usage-thumb.png)](http://vimeo.com/user3503356/godi-usage)

## Development Status

[![Build Status](https://travis-ci.org/Byron/godi.svg?branch=master)](https://travis-ci.org/Byron/godi)

Latest releases [can be found here](https://github.com/Byron/godi/releases).
Developer docs can be found at [godoc.org](http://godoc.org/github.com/Byron/godi)

## Infrastructure

* [Manual and User Documentation](http://byron.github.io/godi)
* [Issue Tracker](https://github.com/Byron/godi/issues)
* [Support Forum](https://groups.google.com/forum/#!forum/go-data-integrity)
* [API Documentation](http://godoc.org/github.com/Byron/godi)
* [Continuous Integration](https://travis-ci.org/Byron/godi)

## Credits

`godi` uses the following libraries:

* [codegangsta/cli](https://github.com/codegangsta/cli)

`godi` is inspired by

* [media hash list](http://mediahashlist.org)

### LICENSE

This open source software is licensed under [GNU Lesser General Public License](https://github.com/Byron/godi/blob/master/LICENSE.md)
