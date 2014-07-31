There are many programs providing a similar feature-set, as hashing data is not exactly a new invention.
This comparison will therefore focus on the tool that inspired the creation of `godi` in the first place, `mhl`.

As the latter doesn't provide the ability to copy data while it hashes it, a program called *DoubleData* or *ShotPutPro* could be used for comparison in that regard. However, as there is no free version of these available, there is no comparison just yet.


## Comparison to MHL

*MHL*, [media-hash-list](http://mediahashlist.org), is a command-line utility to seal and verify data. It's somewhat specialized to data produced in the media industry.

### Benefits over MHL

* **Performance**
    + `godi` is up to [6x faster](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/mhl_performance-comparison.mov.gif) in tests on a quad-core laptop CPU, but will scale with more cores and higher available device bandwidth to reach even higher speeds.
        + Those inclined may maximize bandwidth by [tuning parallelism](details.md#Performance Considerations)
    + Will read in parallel from multiple devices, leveraging all device's maximum bandwidth. 
    + Can write to multiple devices at once, creating multiple duplicates of your data as fast as your slowest writer. If one device [fails](details.md#Error Handling), all other devices will still receive their duplicates, whereas the failed device will [not have an intermediate result](details.md#Atomic Operation).
* **Usability**
    + `godi` just works. `mhl` will fail (for some reason) if it finds files it doesn't want to handle, and it ignores symbolic links. `godi` has no implicit ignore lists, and may be fully customized to suit your needs. By default, it will just ignore volatile files, like *.DS_Store*.
    + `mhl` has a gather stage to investigate all files before it does anything. Even though this information is used to show how much more data it has to process, it will take a lot of time on big datasets. `godi` will start in the moment you hit enter, and inform you about any progress regularly.
        - For some reason, `mhl` presents the total amount of bytes to process according to the amount of bytes it has to hash. [This can be confusing](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/mhl_metrics.mov.gif).
    + `godi` comes with a state of the art commandline interface, allowing to learn the command by using it. No manual required. `mhl` has it as well, but it's inconsistent, the front page (`mhl [help]`) doesn't inform about all available sub-commands and help topics. You will find the missing ones referenced in help-texts of the available topics.
    + `mhl` seals are not protected against being altered. File corruption or intentional changes can't be detected, and will lead to invalid verification results. 
    + `mhl` doesn't escape characters interpreted by XML, producing invalid MHL files in case a `&` character is in a file path, without telling you. This can be a real issue, have a look at [this demonstration to see it for yourself](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/gif/mhl_invalid-seal.mov.gif)
        - I would have left a note, but the [discussion forum](http://mediahashlist.org/discussion/) doesn't exist anymore.
    + During verify, `mhl` will just check one of the sha1 or md5 hashes, in that order, which leaves room as tiny as a billionth atom to accidentally verify a file. `godi` will verify against all hashes it has taken.
    + `godi` is available for more platforms and architectures, `mhl` is 'only' available for windows, linux and osx (all 64bit).
* **Features**
    * *Copy or archive on the fly*
        + While hashing, you can also transfer the data, reading it only once in the process. With `mhl`, you need to copy first, and hash afterwards, which reads the data twice. `godi`s operation assumes the storage works correctly, however, there is a safe mode which verifies the copy nonetheless.
        + It will never overwrite existing files.
    * *Atomic Operation*
        + It will not produce intermediate results, and either finish successfully, or leave no intermediate files at all
        + Particularly useful when copying or archiving, as it will not leave any written file(s), allowing to safely abort and retry at will. The latter is good during performance tuning.
* **Open-Source and active Development**
    + Everyone is free to look at and modify `godi`s source code. `mhl` is closed source, downloads are behind a form to retrieve your e-mail address, and the last life-sign was posted [July 18th 2013](http://mediahashlist.org/news/)

### Benefits of MHL over godi

* **Executable Size and Memory Consumption**
    + The `mhl` executable is only *131kB* (OSX) to *1.5MB* in size, whereas `godi` clocks in at *3.5mB*. This is due to differences in how the program is compiled.
    + `mhl` just uses about 712kB of memory when processing small amounts of files, whereas `godi` needs about 4mB for the same dataset.
* **Automatic Handling of Meta-Data**
    + `mhl` will include more information, including a full operation log, in the generated *mhl file*, whereas `godi` leaves gathering this information to the user.
* **Special Handling for Image Sequences**
    + As `mhl` comes from the media industry, it may read and check frame sequences specifically.
* **Advanced Seal File Handling**
    + `mhl` is able to produce hashes individually and create mhl files from the output of various commands. For example `mhl hash -c /path/to/folder/*.mov | mhl file -s -v -o /path/to/folder` and `openssl dgst -md5 /path/to/files/ -name "*.mov" | mhl file -s -v -o /path/to/folder` are possible applications.


