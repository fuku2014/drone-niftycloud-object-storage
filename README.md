# drone-niftycloud-object-storage

[![Go Doc](https://godoc.org/github.com/fuku2014/drone-niftycloud-object-storage?status.svg)](http://godoc.org/github.com/fuku2014/drone-niftycloud-object-storage)
[![Go Report](https://goreportcard.com/badge/github.com/fuku2014/drone-niftycloud-object-storage)](https://goreportcard.com/report/github.com/fuku2014/drone-niftycloud-object-storage)

Drone plugin to publish files and artifacts to NIFTY Cloud Object Storage. For the
usage information and a listing of the available options please take a look at
[the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/niftycloud-object-storage .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-niftycloud-object-storage' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e PLUGIN_SOURCE=<source> \
  -e PLUGIN_TARGET=<target> \
  -e PLUGIN_BUCKET=<bucket> \
  -e NIFTY_ACCESS_KEY_ID=<token> \
  -e NIFTY_SECRET_KEY=<secret> \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/niftycloud-object-storage
```
