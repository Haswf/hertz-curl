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

var requestToCurlFunc = []struct {
	name string
	f    func(req *protocol.Request) string
}{
	{
		name: "hertz-curl",
		f: func(req *protocol.Request) string {
			cmd := GetCurlCommand(req)
			return cmd.String()
		},
	},
	{
		name: "adaptor and http2curl",
		f: func(req *protocol.Request) string {
			r, _ := adaptor.GetCompatRequest(req)
			// log curl command
			cmd, _ := http2curl.GetCurlCommand(r)
			return cmd.String()
		},
	},
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

func BenchmarkHertzRequest2Curl(b *testing.B) {
	for _, v := range bodyLength {
		for _, impl := range requestToCurlFunc {
			b.Run(fmt.Sprintf("body_length_%d_%s", v.input, impl.name), func(b *testing.B) {
				// run the Fib function b.N times
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPost)
				body := make([]byte, v.input)
				rand.Read(body)
				req.Header.Set(consts.HeaderContentType, "application/json")
				req.SetBody(body)
				req.SetRequestURI("https://www.cloudwego.io/zh/")

				for i := 0; i < b.N; i++ {
					impl.f(req)
				}
			})
		}

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
