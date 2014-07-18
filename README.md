`godi` stands for "go data integrity" and is a commandline utility to generate signature files from a directory tree. This allows to re-check the tree for consistency, and thus verify the data is intact. This is especially useful if data is retrieved from unreliable media, and copied to another storage device.

As it is very common to verify copy operations, `godi` is able to copy files in the moment is hashes them, optionally verifying the destination after it was copied.

## Usage

```bash
# Generate a signature for all files in directory tree/
godi seal tree/

# results in godi-seal.xml file
```


TODO: 

* nprocs - specify how many parallel gather routines there are
* abort-on-error - if False, we continue as long as possible, otherwise we abort and interrupt all currently running procedures
* log-mode - either off, or verbose, and in future, maybe even a binary one which provides a whole lot of additional information

## Performance

Intermediate results indicate a throughput of up to 900MB/s on 2 cores, which is a little more than twice as fast as the single-threaded [mhl](http://mediahashlist.org/).

I am still wondering why it doesn't benefit from more cores.

```bash
$ time  ./godi seal ~/Movies
Sealed 479 files with total size of 407.74786MB in 0.47895445400000003s (851.3290989659139 MB/s, 0 errors)

real    0m0.486s
user    0m0.879s
sys 0m0.076s
```

```bash
$ time mhl seal -v -t sha1 -o ~/Movies  ~/Movies
----------------------
Finished generating checksums for: 
   480 file(s) 
   with total filesize of 407 MB (427719586 bytes)
----------------------
Summary:
   480 of 480 files SUCCEEDED
-------------------
End of input.
Finish date in UTC: 2014-07-17 22:02:42
MHL file path(s):
   /Users/byron/Movies/Movies_2014-07-17_220241.mhl
===================

real    0m1.186s
user    0m1.100s
sys 0m0.085s
```

## Development Status

[![Build Status](https://travis-ci.org/Byron/godi.svg?branch=master)](https://travis-ci.org/Byron/godi)
![under construction](https://raw.githubusercontent.com/Byron/bcore/master/src/images/wip.png)

### LICENSE

This open source software is licensed under [GNU Lesser General Public License](https://github.com/Byron/godi/blob/master/LICENSE.md)
