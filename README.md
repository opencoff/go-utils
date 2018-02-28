# go-utils - general purpose helper functions for golang

## What is it?
This is a collection of library functions I used in multiple go projects.
It didn't have a good home so, I collectively put it here.

## What is available?

 - Threadsafe fixed-size circular queue
 - Random UUIDv4 generator
 - mmap(2) reader to read and process very large files in chunks
 - Channel backed, fixed-size buffer pool. Unlike sync.Pool, this has
   a fixed size (set at construction time) and never changes. As a result,
   when the pool runs out of memory, the caller is blocked until another
   go-routine frees a buffer.
 - Interactive password prompter.


