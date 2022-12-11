[![Build Status](https://travis-ci.org/KarpelesLab/vfs.svg)](https://travis-ci.org/KarpelesLab/vfs)
[![GoDoc](https://godoc.org/github.com/KarpelesLab/vfs?status.svg)](https://godoc.org/github.com/KarpelesLab/vfs)
[![Coverage Status](https://coveralls.io/repos/github/KarpelesLab/vfs/badge.svg?branch=master)](https://coveralls.io/github/KarpelesLab/vfs?branch=master)

# Filesystem Abstraction in Go

Yet another one, created because none [of the bazillon existing Golang VFS](https://awesome-go.com/#files) matched
the needs we have.

Specifically, many tend to be either oriented toward local filesystems or cloud
storage providers. Typically, both are very different. Local filesystems allow
files to be modified, while cloud providers typically require a whole file to
be re-uploaded for any change. As such, cloud-oriented libraries may support
local filesystem but have no API for partial writes locally.

Here, cloud storage solutions such as AWS S3 are considered "keyvals", similar
to databases where changing a byte in a value requires rewriting the whole
value.
The goal is to be able to offer converter interfaces that expose such
backends as proper filesystems supporting partial writes.

# Focus

This implementation focuses on the following goals:

* Stay as close as possible to filesystem concepts
* Be as compatible as possible with Golang's interfaces
* Be as simple as possible to extend
* Allow working with limited key/value backends

# Features

* Filesystem Backends:
  * localfs filesystem
  * memfs
  * zipfs (read only)
* Keyval Backends:
  * memkv
  * boltkv using [boltdb](https://github.com/boltdb/bolt)
* Converters:
  * vdirfs: provides directory indexation/listing for backends which do not have this feature (such as zipfs)

## Planned

* Support for a wide range of backends (AWS S3, etc)
* Support for frontends (fuse, http, etc)
* Middlewares (keyval→filesystem adapters, encryption, etc)
