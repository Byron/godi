![under construction](https://raw.githubusercontent.com/Byron/bcore/master/src/images/wip.png)

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


## TODO: About [DoubleData](http://www.doubledata.biz)