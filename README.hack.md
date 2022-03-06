# Docs on how to hack here!

**TL;DR** just `go build -v`

## 101s

Install `Go`: https://go.dev/dl make sure to do the `export PATH=$PATH:$(go env GOPATH)/bin`

Generate Prorobuf stubs if API updated: Install `protoc` from https://github.com/protocolbuffers/protobuf/releases/ then run `go generate -v ./...`

## Run Tests

```shell
go test -v ./...
```

Integration tests and `etcdstore` tests needs `etcd` in path. Install etcd: https://etcd.io/docs/v3.5/install/

## TODO

- add protobuf APIs
- add buf http://buf.build/
