# hertz-curl


Convert Hertz's protocol.Request to CURL command line, especially handy for logging request in a compat and universal format.

This is a community driven project, taken inspiration from [http2curl](https://github.com/moul/http2curl). Thanks [http2curl](https://github.com/moul/http2curl) developers for previous work.



## Example 
```go
package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/haswf/hertz-curl"
)

func main() {
	h := server.Default()

	h.GET("/curl", func(ctx context.Context, c *app.RequestContext) {
		cmd, _ := curl.GetCurlCommand(&c.Request)
		hlog.Info(cmd.String())
		// 2023/04/06 22:33:00.044947 main.go:15: [Info] curl -X 'GET' -H 'Accept: */*' -H 'Accept-Encoding: gzip, deflate, br' -H 'Connection: keep-alive' -H 'Host: localhost:8888' -H 'Postman-Token: bc98e52c-e9fd-4c71-895b-9a27d940f151' -H 'User-Agent: PostmanRuntime/7.29.2' 'http://localhost:8888/curl'
		c.JSON(consts.StatusOK, utils.H{"curl": cmd.String()})
	})

	h.Spin()
}

```

## Install

```bash
go get github.com/haswf/hertz-curl
```

## Why hertz-curl?
1. It is possible to convert protocol.Request to http.Request first with [adaptor](https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/adaptor/), then generate curl command using [http2curl](https://github.com/moul/http2curl).
However, this conversion is costly because req.Body will be copied twice during this process.
2. Both `adaptor.GetCompatRequest` and `http2curl.GetCurlCommand` return an error if:
- http method is invalid or url is malformed (See http/request.go:856). Even if an url is malformed, we may want to get its curl representation for debugging.
- io error when reading req.Body (See http2curl.go:50). This case could never happen in practice. 

In summary, hertz-curl simplify curl generation and significantly improve performance by avoid unnecessary copying.

## Benchmark
```bash
goos: darwin
goarch: amd64
pkg: github.com/haswf/hertz-curl
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkHertzRequest2Curl/body_length_1024_hertz-curl-12                         394759              2879 ns/op
BenchmarkHertzRequest2Curl/body_length_1024_adaptor_and_http2curl-12              201104              5703 ns/op
BenchmarkHertzRequest2Curl/body_length_2048_hertz-curl-12                         335019              3468 ns/op
BenchmarkHertzRequest2Curl/body_length_2048_adaptor_and_http2curl-12              146192              8153 ns/op
BenchmarkHertzRequest2Curl/body_length_4096_hertz-curl-12                         221118              5465 ns/op
BenchmarkHertzRequest2Curl/body_length_4096_adaptor_and_http2curl-12               92596             13106 ns/op
BenchmarkHertzRequest2Curl/body_length_8192_hertz-curl-12                         139407              8292 ns/op
BenchmarkHertzRequest2Curl/body_length_8192_adaptor_and_http2curl-12               56043             21170 ns/op
BenchmarkHertzRequest2Curl/body_length_16384_hertz-curl-12                         77752             14924 ns/op
BenchmarkHertzRequest2Curl/body_length_16384_adaptor_and_http2curl-12              28854             41775 ns/op
BenchmarkHertzRequest2Curl/body_length_32768_hertz-curl-12                         35098             34517 ns/op
BenchmarkHertzRequest2Curl/body_length_32768_adaptor_and_http2curl-12              14514             82818 ns/op
BenchmarkHertzRequest2Curl/body_length_65536_hertz-curl-12                         19610             58740 ns/op
BenchmarkHertzRequest2Curl/body_length_65536_adaptor_and_http2curl-12               7759            148979 ns/op
```

## License
This project is under the Apache License 2.0. See the LICENSE file for the full license text.