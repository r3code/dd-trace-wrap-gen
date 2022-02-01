dd-trace-wrap-gen
========

Generates interface decorators with [DataDog Distributes Tracing](https://github.com/DataDog/dd-trace-go) support.

Installation
------------

```
go get github.com/r3code/dd-trace-wrap-gen
```

Example
-------

```go
type Service interface {
	Set(ctx context.Context, key, value []byte) error
	Get(ctx context.Context, key []byte) (value []byte, err error)
}
```

Usage
-----

```
dd-trace-wrap-gen -i <an interface to decortate> -o <output_file> ./example
```

* -i - specify interface name  
* -o - output filename  
* -s - target struct name, default: <interface name>WithTracing  


```
dd-trace-wrap-gen -i Service -o example/service_trace.go ./example
```

Will generate:
```go

PUT EXAMPLE HERE
```
