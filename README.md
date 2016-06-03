# Buster [![GoDoc](https://godoc.org/github.com/james-maloney/buster?status.svg)](https://godoc.org/github.com/james-maloney/buster)

A Go file server used to bust static file caches on startup. This lets us set a very aggresive, 1 year cache control header
and ensures that the client gets fresh files on starup since the file path changes. If the client does request using an old file path
the new file will be served.

## Usage

See the example [code](https://github.com/james-maloney/buster/blob/master/example/example.go) for basic usage.
