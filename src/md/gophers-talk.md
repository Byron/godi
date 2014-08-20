## About the Speaker

Sebastian is a software engineer rooted in the VFX industry. Having worked on technical challenges in movies like "Iron Man 3" and "Captain America: The Winter Solder", he took care of automating IT systems and offloading expensive tasks to the compute farm. 'Written languages' are python, c++ and by now, go.

## Abstract about: Lessons learned when making 'godi'

Learning a new programming language is best achieved by solving an actual problem. `godi` is the project I used to learn `go` and I will share the most exciting things about go that were revealed to me in the process. The speech will be held from the perspective of a go-beginner, showing aspects of the go language using actual code of `godi`.

Among these aspects are:

* Motivation: why use go ?
* Why godi ?
    + features
    + architecture
* Typical go project structure and project location
    + what to know about the main package
* About go executables (size on disk, memory consumption, dependencies)
    + ... and how to reduce their size
    + cross-platform development and platform independent code
    + build tags for platform specific code
* Test-driven development
* Performance
    + achieving parallelism and troubles on the way
    + share state by communicating
        - why a shared state can be useful too
    + about heap pressure and interfaces
* Error Handling
    + handle each error
    + communication from consumers to producers
* Code documentation
* Source and binary distribution