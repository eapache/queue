Queue
=====

A fast Golang queue using a ring-buffer, based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.

The queue implemented here is as fast as it is in part because it is *not* thread-safe.
