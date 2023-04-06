package curl

import (
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"net/url"
	"testing"
)

func TestGetCurlCommand(t *testing.T) {
	doc := "https://cloudwego.github.io"
	tests := []struct {
		Name    string
		Request func() *protocol.Request
		Want    string
	}{
		{
			Name: "simple get request",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodGet)
				req.SetRequestURI(doc)
				return req
			},
			Want: "curl -X 'GET' 'https://cloudwego.github.io/'",
		},
		{
			Name: "get request with query",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodGet)
				query := url.Values{}
				query.Set("key", "value")
				req.SetRequestURI(doc)
				req.SetQueryString(query.Encode())
				return req
			},
			Want: "curl -X 'GET' 'https://cloudwego.github.io/?key=value'",
		},
		{
			Name: "post request with query, header and body",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPost)
				query := url.Values{}
				query.Set("key", "value")
				req.SetRequestURI(doc)
				req.SetQueryString(query.Encode())
				req.SetHeader(consts.HeaderContentType, "application/json")
				req.SetBody([]byte(`{"hello": "hertz"}`))
				return req
			},
			Want: `curl -X 'POST' -d '{"hello": "hertz"}' -H 'Content-Type: application/json' 'https://cloudwego.github.io/?key=value'`,
		},
	}

	for _, tt := range tests {
		if got := GetCurlCommand(tt.Request()).String(); got != tt.Want {
			t.Fatalf("%s error: want %s, got %s", tt.Name, tt.Want, got)
		}
	}
}
