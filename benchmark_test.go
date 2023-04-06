package curl

import (
	"fmt"
	"github.com/cloudwego/hertz/cmd/hz/util"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"math/rand"
	"moul.io/http2curl"
	"testing"
)

var bodyLength = []struct {
	input int
}{
	{input: 1024},
	{input: 2048},
	{input: 4096},
	{input: 8192},
	{input: 16384},
	{input: 32768},
	{input: 65536},
}

var byte2StrFunc = []struct {
	name string
	impl func(in []byte) string
}{
	{
		name: "util.Bytes2Str",
		impl: util.Bytes2Str,
	},
	{
		name: "string cast",
		impl: func(in []byte) string {
			return string(in)
		},
	},
}

func BenchmarkHTTP2Curl(b *testing.B) {
	for _, v := range bodyLength {
		b.Run(fmt.Sprintf("body_length_%d", v.input), func(b *testing.B) {
			// run the Fib function b.N times
			req := &protocol.Request{}
			req.SetMethod(consts.MethodPost)
			body := make([]byte, v.input)
			rand.Read(body)
			req.Header.Set(consts.HeaderContentType, "application/json")
			req.SetBody(body)
			req.SetRequestURI("https://www.cloudwego.io/zh/")

			for i := 0; i < b.N; i++ {
				r, _ := adaptor.GetCompatRequest(req)
				// log curl command
				cmd, _ := http2curl.GetCurlCommand(r)
				_ = cmd.String()
			}
		})
	}

}

func BenchmarkHertzCurl(b *testing.B) {

	for _, v := range bodyLength {
		b.Run(fmt.Sprintf("body_length_%d_", v.input), func(b *testing.B) {
			// run the Fib function b.N times
			req := &protocol.Request{}
			req.SetMethod(consts.MethodPost)
			body := make([]byte, v.input)
			rand.Read(body)
			req.Header.Set(consts.HeaderContentType, "application/json")
			req.SetBody(body)
			req.SetRequestURI("https://www.cloudwego.io/zh/")

			for i := 0; i < b.N; i++ {
				cmd := GetCurlCommand(req)
				_ = cmd.String()

			}
		})
	}

}

func BenchmarkByte2Str(b *testing.B) {
	for _, v := range bodyLength {
		for _, f := range byte2StrFunc {
			b.Run(fmt.Sprintf("body_length_%d_func_%s", v.input, f.name), func(b *testing.B) {
				// run the Fib function b.N times
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPost)
				body := make([]byte, v.input)
				rand.Read(body)
				req.Header.Set(consts.HeaderContentType, "application/json")
				req.SetBody(body)
				req.SetRequestURI("https://www.cloudwego.io/zh/")

				for i := 0; i < b.N; i++ {
					cmd := parse(req, f.impl)
					_ = cmd.String()

				}
			})
		}

	}

}
