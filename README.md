# simian

[![GoDoc](https://godoc.org/github.com/mandykoh/simian?status.svg)](https://godoc.org/github.com/mandykoh/simian)

Caution—This code is proof-of-concept quality. This means:

  * The API is unstable and doesn’t support common use cases.
  * The index format is unstable.
  * It is not thread safe.
  * Test coverage is mostly non-existant.
  * The current fingerprinting method has known weaknesses which affect quality of results.

Development
-----------

Simian uses [govendor][govendor] to manage dependencies.

```bash
go get github.com/mandykoh/simian
cd $GOPATH/src/github.com/mandykoh/simian

# Install dependencies
govendor sync
govendor install +local

# Run the test suite
go test ./... | grep -v /vendor/
```

  [govendor]: https://github.com/kardianos/govendor
