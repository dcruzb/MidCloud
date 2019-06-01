# MidCloud [![Godoc](https://godoc.org/github.com/dcbCIn/MidCloud?status.svg)](https://godoc.org/github.com/dcbCIn/MidCloud)
Adaptive middleware for cloud services, with transparent choice of best cloud.
The basic conceptions are that the message will be sent to the server that has the lower cost to do the task and is available.

## Installation

Standard `go get`:

```
$ go get github.com/dcbCIn/MidCloud
```

## Usage & Example

For usage see the [MidCloud Godoc](http://godoc.org/github.com/dcbCIn/MidCloud).

For examples see [CloudStorage Godoc](http://godoc.org/github.com/dcbCIn/CloudStorage).
The [CloudStorage](https://github.com/dcbCIn/CloudStorage) project has examples associated with it there.

## But Why?!

There exists a long list of middlewares, but MidCloud comes with the proposal of reduce the costs of multi-cloud 
applications, and to eradicate any downtime possible. 

Can be from a multi-cloud storage (see [CloudStorage](https://github.com/dcbCIn/CloudStorage) project), multi-cloud faas, 
or any other you want to create).