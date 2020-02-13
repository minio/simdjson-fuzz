# simdjson-fuzz

Fuzzers and corpus for [github.com/minio/simdjson-go](https://github.com/minio/simdjson-go)

# Running

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

Feel free to submit additional corpus as Pull Requests.

# Timeouts

Due to problems with the current Go fuzzer better results can often be obtained by running the fuzzer for
a shorter time period and restarting it.

A Go program is supplied to do that will run the supplied command for a specific duration.

A continuously looping script can be created in bash for instance:

```
go-fuzz-build -o=fuzz-build.zip -func=FuzzCorrect .
while true; do
    go run timeout.go -duration=10m go-fuzz -bin=fuzz-build.zip -workdir=corpus`
done
```

In the example the fuzzer runs for 10 minutes before exiting.
