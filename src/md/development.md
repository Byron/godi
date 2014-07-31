
## Package Overview

`godi` is a small package with just a few parts, as illustrated by the following diagram.

![architecture](https://raw.githubusercontent.com/Byron/godi/web-resources/lib/png/arch.txt.png)

As you can see, the *api*, *utility* and *io* packages are rather standalone and don't require much to work, whereas the highest level packge dealing with the command-line interface, *cli*, directly and indirectly uses all packages underneath.

