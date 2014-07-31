/*
Package io provides asynchronous reader and writer implementations.

Next to hashing, the bulk of the time is spent doing input and output operations.
For devices to achieve their maximum bandwidth, it's important not to stress them
with too many parallel IOPs, which makes it necessary to control the amount of
input/output streams. This is the key-feature of the ReadChannelController and
WriteChannelController types implemented here.

Similar to io.MultiWriter, there is a ParallelMultiWriter implementation which sends
bytes to multiple writers at once.
*/
package io
