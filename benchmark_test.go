package curl

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"moul.io/http2curl"
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
	f    func(req *protocol.Request) (string, error)
}{
	{
		name: "hertz-curl",
		f: func(req *protocol.Request) (string, error) {
			cmd, _ := GetCurlCommand(req)
			return cmd.String(), nil
		},
	},
	{
		name: "adaptor and http2curl",
		f: func(req *protocol.Request) (string, error) {
			r, err := adaptor.GetCompatRequest(req)
			if err != nil {
				return "", err
			}
			// log curl command
			cmd, err := http2curl.GetCurlCommand(r)
			if err != nil {
				return "", err
			}
			return cmd.String(), nil
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
					_, _ = impl.f(req)
				}
			})
		}
	}
}
