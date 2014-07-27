`godi` stands for "go data integrity" and is a commandline utility to generate signature files from a directory tree. This allows to re-check the tree for consistency, and thus verify the data is intact. This is especially useful if data is retrieved from unreliable media, and copied to another storage device.

As it is very common to verify copy operations, `godi` is able to copy files in the moment is hashes them, optionally verifying the destination after it was copied.

## Usage

```bash
# Generate a signature for all files in directory tree/
godi seal tree/

# results in godi-seal.xml file
```

## Benefits over MHL

* **Performance**
    + `godi` is up to multiple times faster
    + Those inclined may maximize bandwidth by tuning parallelism
    + Will read in parallel from multiple devices, leveraging all device's maximum bandwidth.
    + Can write to multiple devices at once, creating multiple duplicates of your data as fast as your slowest writer. If one device fails, all other devices will still receive their duplicates, whereas the failed device will not have an intermediate result (see atomic operation).
* **Copy or archive on the fly**
    + While hashing, you can also transfer the data, reading it only once in the process. With MHL, you need to copy first, and hash afterwards, which reads the data twice. `godi`s operation assumes the storage works correctly, however, there is a safe mode which verifies the copy nonetheless.
    + It will never overwrite existing files.
* **Atomic Operation**
    + It will not produce intermediate results, and either finish successfully, or not at all.
    + Particularly useful when copying or archiving, as it will not leave any written file(s), allowing to safely abort and retry at will. The latter is good during performance tuning.
* **Usability**
    + `godi` just works. `mhl` will fail (for some reason) if it finds a hidden file. `godi` will just ignore hidden files and symbolic links and otherwise process everything in its way.
    + `godi` comes with a state of the art commandline interface, allowing to learn the command by using it. No manual required.
    + `mhl` seals are not protected against being altered. File corruption or intentional changes can't be detected, and will lead to invalid verification results.
    + During verify, `mhl` will just check one of the sha1 or md5 hashes, in that order, which leaves room as tiny as a billionth atom to accidentally verify a file. `godi` will verify against all hashes it has taken.


## Limitations

### Windows
* Multi-device optimizations [don't currently apply](https://github.com/Byron/godi/issues/13) on windows
* When ctrl+C is pressed in the in the git-bash to interrupt the program, godi will attempt to stop, but appears to be killed before it can finish cleanup. This seriously hampers atomic operation, and it is advised to use the cmd.exe prompt. Might be related to [this issue](http://stackoverflow.com/questions/10021373/what-is-the-windows-equivalent-of-process-onsigint-in-node-js) in some way.

### General
* Sealed copy ignores permission bits on directories, and will create them with `0777` in generally. It does, however, respect and maintain the mode of copied files.
* `godi` is very careful about memory consumption, yet atomicity comes at the cost of keeping a list of files already copied for undo purposes. That list grows over time, and consumed ~200MB for 765895 files. It might be worth providing a flag to turn undo off.

## Development Status

[![Build Status](https://travis-ci.org/Byron/godi.svg?branch=master)](https://travis-ci.org/Byron/godi)
![under construction](https://raw.githubusercontent.com/Byron/bcore/master/src/images/wip.png)

## Credits

`godi` uses the following libraries:

* [codegangsta/cli](https://github.com/codegangsta/cli)

### LICENSE

This open source software is licensed under [GNU Lesser General Public License](https://github.com/Byron/godi/blob/master/LICENSE.md)
