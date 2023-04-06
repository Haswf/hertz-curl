package curl

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func TestGetCurlCommand(t *testing.T) {
	uri := "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu"
	tests := []struct {
		Name    string
		Request func() *protocol.Request
		Want    string
	}{
		{
			Name: "get request",
			Request: func() *protocol.Request {
				form := url.Values{}
				form.Add("age", "10")
				form.Add("name", "Hudson")
				body := form.Encode()

				req := &protocol.Request{}
				req.SetMethod(consts.MethodPost)
				req.SetRequestURI("http://foo.com/cats")
				req.Header.Set("API_KEY", "123")
				req.SetBody([]byte(body))
				return req
			},
			Want: "curl -X 'POST' -d 'age=10&name=Hudson' -H 'Api_key: 123' 'http://foo.com/cats' --compressed",
		},
		{
			Name: "json",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPut)
				req.SetRequestURI(uri)
				req.Header.Set("Content-Type", "application/json")
				req.SetBody([]byte(`{"hello":"world","answer":42}`))

				return req
			},
			Want: "curl -X 'PUT' -d '{\"hello\":\"world\",\"answer\":42}' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu' --compressed",
		},
		{
			Name: "no body",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPut)
				req.Header.Set("Content-Type", "application/json")
				req.SetRequestURI(uri)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			Want: "curl -X 'PUT' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu' --compressed",
		},
		{
			Name: "empty string body",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPut)
				req.Header.Set("Content-Type", "application/json")
				req.SetBody([]byte(""))
				req.SetRequestURI(uri)
				return req
			},
			Want: "curl -X 'PUT' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu' --compressed",
		},
		{
			Name: "new line in body",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPost)
				req.Header.Set("Content-Type", "application/json")
				req.SetBody([]byte("hello\nworld"))
				req.SetRequestURI(uri)
				return req
			},
			Want: "curl -X 'POST' -d 'hello\nworld' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu' --compressed",
		},
		{
			Name: "special characters in body",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPost)
				req.Header.Set("Content-Type", "application/json")
				req.SetBody([]byte(`Hello $123 o'neill -"-`))
				req.SetRequestURI(uri)
				return req
			},
			Want: "curl -X 'POST' -d 'Hello $123 o'\\''neill -\"-' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu' --compressed",
		},
		{
			Name: "combined",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPut)
				req.SetRequestURI(uri)
				req.Header.Set("X-Auth-Token", "private-token")
				req.Header.Set("Content-Type", "application/json")
				req.SetBody([]byte(`{"hello":"world","answer":42}`))

				return req
			},
			Want: `curl -X 'PUT' -d '{"hello":"world","answer":42}' -H 'Content-Type: application/json' -H 'X-Auth-Token: private-token' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu' --compressed`,
		},
		{
			Name: "https",
			Request: func() *protocol.Request {
				req := &protocol.Request{}
				req.SetMethod(consts.MethodPut)
				req.SetRequestURI("https://www.example.com/abc/def.ghi?jlk=mno&pqr=stu")
				req.Header.Set("X-Auth-Token", "private-token")
				req.Header.Set("Content-Type", "application/json")
				req.SetBody([]byte(`{"hello":"world","answer":42}`))

				return req
			},
			Want: `curl -k -X 'PUT' -d '{"hello":"world","answer":42}' -H 'Content-Type: application/json' -H 'X-Auth-Token: private-token' 'https://www.example.com/abc/def.ghi?jlk=mno&pqr=stu' --compressed`,
		},
	}

	for _, tt := range tests {
		cmd := GetCurlCommand(tt.Request())
		if cmd.String() != tt.Want {
			t.Fatalf("%s error: want %s, got %s", tt.Name, tt.Want, cmd.String())
		}
	}
}

func TestGetCurlCommand_ServerAndClientSide(t *testing.T) {
	h := server.Default()
	var serverCurl *Command
	// run hertz server in a new goroutine as h.Spin blocks.
	go func() {
		h.GET("/curl", func(ctx context.Context, c *app.RequestContext) {
			serverCurl = GetCurlCommand(&c.Request)
			c.JSON(consts.StatusOK, nil)
		})
		h.Spin()
	}()

	// wait for the server to run
	for !h.IsRunning() {
		time.Sleep(time.Millisecond)
	}

	c, err := client.NewClient()
	if err != nil {
		return
	}
	var resp struct {
		Curl string `json:"curl"`
	}

	req := &protocol.Request{}
	res := &protocol.Response{}
	req.SetMethod(consts.MethodGet)
	req.SetRequestURI("http://127.0.0.1:8888/curl")
	err = c.Do(context.Background(), req, res)
	if err != nil {
		return
	}
	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		return
	}
	clientCurl := GetCurlCommand(req)
	want := `curl -X 'GET' -H 'Host: 127.0.0.1:8888' -H 'User-Agent: hertz' 'http://127.0.0.1:8888/curl' --compressed`
	if serverCurl.String() != want {
		t.Fatalf("server curl error: want %s, got %s", want, serverCurl.String())
	}
	if clientCurl.String() != want {
		t.Fatalf("client curl error: want %s, got %s", want, clientCurl.String())
	}
}
