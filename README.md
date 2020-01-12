[![Build Status](https://travis-ci.org/KarpelesLab/vfs.svg)](https://travis-ci.org/KarpelesLab/vfs)
[![GoDoc](https://godoc.org/github.com/KarpelesLab/vfs?status.svg)](https://godoc.org/github.com/KarpelesLab/vfs)
[![Coverage Status](https://coveralls.io/repos/github/KarpelesLab/vfs/badge.svg?branch=master)](https://coveralls.io/github/KarpelesLab/vfs?branch=master)

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
