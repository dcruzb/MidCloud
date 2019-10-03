# MidCloud [![Godoc](https://godoc.org/github.com/dcbCIn/MidCloud?status.svg)](https://godoc.org/github.com/dcbCIn/MidCloud) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dcbCIn/MidCloud/blob/master/LICENSE) [![Build Status](https://travis-ci.org/dcbCIn/MidCloud.png?branch=master)](https://travis-ci.org/dcbCIn/MidCloud)
Adaptive middleware for cloud services, with transparent choice of best cloud.
The basic conceptions are that the message will be sent to the server that has the lower cost to do the task and is available.

## Installation

Standard `go get`:

```
$ go get -u github.com/dcbCIn/MidCloud/...
```

In your source code:

```go
import (
	"github.com/dcbCIn/MidCloud/distribution"
	"github.com/dcbCIn/MidCloud/services/common"
)
```

## Usage & Example

For usage see the [MidCloud Godoc](http://godoc.org/github.com/dcbCIn/MidCloud).

For examples see [CloudStorage Godoc](http://godoc.org/github.com/dcbCIn/CloudStorage).
The [CloudStorage](https://github.com/dcbCIn/CloudStorage) project has examples associated with it there.

## But Why?!

There exists a long list of adaptive middlewares, but MidCloud comes with the proposal of reduce the costs of multi-cloud 
applications, and to eradicate any downtime possible. 

Can be from a multi-cloud storage (see [CloudStorage](https://github.com/dcbCIn/CloudStorage) project), multi-cloud faas, 
or any other service you want to create).

## License

[The MIT License (MIT)](https://github.com/dcbCIn/MidCloud/blob/master/LICENSE)
