This document aims at users with at least basic knowledge on how to use a terminal to run command-line programs.

## How to Run Godi

Godi is a command-line program, which requires you to use a terminal emulator to run it and see its output. If you are unfamiliar with terminals or don't have the time to look into how to use one, `godi` isn't for you.

You can at all times and without negative effects [abort a running command](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_sealed-copy_cancelled.mov.gif) by pressing `ctrl+c`, which is the same as for all other command-line tools.

## Helping Yourself

It's not always feasible to consult the Internet to learn how to do something. Good to know that `godi` has all the help you will ever need to use it built right into itself.

It has a modern command-line interface which uses sub-commands to separate distinct features from each other. In any case, you will always see an elaborate help text as follows.

```bash
# General help
$ godi
# The same as above
$ godi help
# Help for the specified sub-command
$ godi help sealed-copy

# Traditional ways to do that work as well
$ godi -h
$ godi --help
# and the same for subcommands
$ godi seal -h
$ godi seal --help
```

## Subcommands

The following sections explain which `godi` sub-command to use based on the particular problem it solves for you.

### Seal - Protect Data from Change
![seal](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_seal.mov.gif)

The *seal* sub-command produces a *seal file* which stores signatures of files. A signature is like a fingerprint, uniquely identifying the entire contents of the file. If that contents changes in any way, it's signature changes with it.

Therefore a seal is like an insurance you want to have when dealing with valuable, immutable data to assure it still is what you think it is.

You can seal one or more files or directories in one go.

```bash
# Seal a file and a directory, producing two seal files
$ godi seal ~/Desktop/myWeddingVideo.mov /Volumes/encrypted/taxes/2012
```

This sub-command is affected by [input file filters](details.md#Input File-Filters), and it is possible to choose between [different seal formats](details.md#Seal Formats). You may also be interested in information about [error handling policy](details.md#Error Handling)

### Verify - Assure Data didn't Change
![seal](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_verify.mov.gif)

This operation requires a *seal file* and will let you know if the files recorded within the seal have changed on disk, have changed their size, or are missing entirely. If you do not have a seal file, you can create one using the *seal* sub-command.

You can specify one or more seal files to be checked at once.

```bash
# Verify files on disk have not been altered compared to the given seal file
$ godi verify ~/Desktop/godi_2014-07-30_102257.gobz /Volumes/encrypted/taxes/2012/godi_2014-07-30_102259.gobz
```

And this is how it looks if [something is not in order](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_verify-fail.mov.gif).

If a source file cannot be read, *verify* will continue with other files to provide as much information to you as possible.

### Sealed Copy - Seal with Duplication
![sealed-copy](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_sealed-copy.mov.gif)

Considering that a seal will only help you to identify undesired changes to your valuable data, the only way to restore unaltered versions of such data is to have multiple copies if it. In information technology, the safety of data is expressed in terms of redundancy, that is how many copies of data is available. There are various ways to increase the redundancy and therefore safety of data, one of which `godi` can help you with.


The *sealed-copy* copy sub-command allows you to not only seal your data, but to generate one or more duplicates in the same go. This is preferred to copy manually, and seal afterwards, as it will read your data only once, and will therefore be faster.

```bash
# Copy to a single destination - godi does not create the destination directory for you
$ godi sealed-copy ~/valuables /Volumes/backup/valuables

# You can create multiple copies at once - note the double-dash to separate sources from destinations
$ godi sealed-copy ~/valuables -- /Volumes/backup01/valuables /Volumes/backup02/valuables
```

If you have the time and if you want to absolutely sure that the destination data is exactly what you expect, you can use the `--verify` flag. That way, all data will be verified right after all copy operations are finished. Please also note that this will re-read all written data, which takes additional time.

```bash
# Will re-read all data in the destination to assure it was written correctly.
$ godi sealed-copy --verify ~/valuables /Volumes/backup/valuables
```

Have a look at [this video](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_sealed-copy-verify_full.mov.gif) to see the `--verify` flag in action.

This sub-command is affected by [input file filters](details.md#Input File-Filters), and subject to [atomic operations](details.md#Atomic Operation). You may also be interested to learn how it deals with [errors](details.md#Error Handling) while writing to a destination.

`godi` will *never* overwrite existing files, as shown [in this video](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/godi_sealed-copy_fail-write.mov.gif).

