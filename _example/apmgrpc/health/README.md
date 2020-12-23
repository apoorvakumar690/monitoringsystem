# DEVELOPMENT

Health monitoring by sending requests to an endpoint on the application.

## Generate gRPC service

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
cd ~/go/src/stash.bms.bz/bms/monitoringsystem/_example/grpc/health
protoc --go_out=plugins=grpc:. health.proto
```

## Installing

```bash
# Installing

# Install via brew (recommended)
brew install protobuf

# Go support for Protocol Buffers - Google's data interchange format
# https://github.com/golang/protobuf
go get -u github.com/golang/protobuf/protoc-gen-go

# Facing any issue in proto-gen-go not found
# I have resolved it like this:

# go get -u github.com/golang/protobuf/protoc-gen-go
# copy protoc-gen-go from bin folder from go workspace to /usr/local/bin/
# note: i've installed grpc and protobuf using homebrew

# protoc need to know where protoc-gen-go is, alternatively (and probably better) is to:

go get -u github.com/golang/protobuf/protoc-gen-go
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

```

### Debug the gRPC proto [grpcurl]

**Installation**

Download the binary from the releases page.

On macOS, grpcurl is available via Homebrew:

```sh
brew install grpcurl
```

[proto3]: https://developers.google.com/protocol-buffers/docs/proto3
[sonar-golang]: https://github.com/uartois/sonar-golang
[grpcurl]: https://github.com/fullstorydev/grpcurl
