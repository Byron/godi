
![under construction](https://raw.githubusercontent.com/Byron/bcore/master/src/images/wip.png)

## Seal Formats

* Also mention performance implications

## Increasing Performance

For understanding this paragraph, it's beneficial to understand how data is processed in godi. Without getting into too much detail, you can see that data is first read from storage, then hashed, and possibly be written in any of the copy-enabled modes.

![architecture](../img/arch.dia.txt.png)

As the *Hasher* part can easily deliver 450MB/s per core, you can imagine that the bottleneck will occour during disk-based input-output operations. For example, reading from an SSD with a cold filesystem cache will rarely deliver more than 500MB/s, and writing to an SSD would not be much faster either.

Nonetheless, depending on the type of storage, you might benefit from multiple simultaneous reads, and/or multiple simultaneous writes, which may drastically increase the perceived performance.

It is vital to test for good values for `--num-readers` and `--num-writers` to get optimal performance for your respective hardware. By default, there may be as many hashers as you have cores, and this rarely needs a change unless `godi` is competing with other programs for the CPU.



## Limitations

### Windows
* Multi-device optimizations [don't currently apply](https://github.com/Byron/godi/issues/13) on windows
* When ctrl+C is pressed in the in the git-bash to interrupt the program, godi will attempt to stop, but appears to be killed before it can finish cleanup. This seriously hampers atomic operation, and it is advised to use the cmd.exe prompt. Might be related to [this issue](http://stackoverflow.com/questions/10021373/what-is-the-windows-equivalent-of-process-onsigint-in-node-js) in some way.

### General
* Sealed copy ignores permission bits on directories, and will create them with `0777` in generally. It does, however, respect and maintain the mode of copied files.
* `godi` is very careful about memory consumption, yet atomicity comes at the cost of keeping a list of files already copied for undo purposes. That list grows over time, and consumed ~200MB for 765895 files. It might be worth providing a flag to turn undo off.

