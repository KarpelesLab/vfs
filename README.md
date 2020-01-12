# Filesystem Abstraction in Go

Yet another one, created because note of the bazillon existing VFS matched
the needs we have.

# Focus

This implementation focuses on the following goals:

* Stay as close as possible to filesystem concepts
* Be as compatible as possible with Golang's interfaces
* Be as simple as possible to extend

# Features (planned)

* Support for a wide range of backends (local, AWS S3, zip file, etc)
* Support for frontends (fuse, http, etc)
