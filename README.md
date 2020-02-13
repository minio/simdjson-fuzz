# simdjson-fuzz
Fuzzers and corpus for [github.com/minio/simdjson-go](https://github.com/minio/simdjson-go)

# running

```
go get -u github.com/minio/simdjson-fuzz
go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build
```

Go to `$GOPATH$/src/github.com/minio/simdjson-fuzz`

Crash testing, execute:
```
go-fuzz-build -o=fuzz-build.zip -func=Fuzz .
go-fuzz -bin=fuzz-build.zip -workdir=corpus
```

Correctness testing, execute:
```
go-fuzz-build -o=fuzz-build.zip -func=FuzzCorrect .
go-fuzz -bin=fuzz-build.zip -workdir=corpus
```

This package does on purpose not use modules.
