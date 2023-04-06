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
