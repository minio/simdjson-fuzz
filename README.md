# simdjson-fuzz
Fuzzers and corpus for [github.com/fwessels/simdjson-go](https://github.com/fwessels/simdjson-go)

# running

```
go get -u github.com/klauspost/simdjson-fuzz
go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build
```

Go to `$GOPATH$/src/github.com/klauspost/simdjson-fuzz`

Execute:
```
go-fuzz-build -o=fuzz-build.zip -func=Fuzz .
go-fuzz -bin=fuzz-build.zip -workdir=corpus
```
