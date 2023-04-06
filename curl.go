package curl

import (
	"fmt"
	"github.com/cloudwego/hertz/cmd/hz/util"
	"github.com/cloudwego/hertz/pkg/protocol"
	"sort"
	"strings"
)

// Command contains exec.Command compatible slice + helpers
type Command struct {
	slice []string
}

// append appends a string to the CurlCommand
func (c *Command) append(newSlice ...string) {
	c.slice = append(c.slice, newSlice...)
}

// String returns a ready to copy/paste command
func (c *Command) String() string {
	return strings.Join(c.slice, " ")
}

func bashEscape(str string) string {
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}

func GetCurlCommand(req *protocol.Request) *Command {
	return parse(req, util.Bytes2Str)
}

// GetCurlCommand returns a CurlCommand corresponding to an http.Request
func parse(req *protocol.Request, byte2str func(in []byte) string) *Command {
	command := Command{}

	command.append("curl")

	command.append("-X", bashEscape(byte2str(req.Method())))

	body := req.Body()
	if len(body) > 0 {
		bodyEscaped := bashEscape(byte2str(body))
		command.append("-d", bodyEscaped)
	}

	var keys []string

	req.Header.VisitAll(func(key, value []byte) {
		keys = append(keys, byte2str(key))
	})
	sort.Strings(keys)

	for _, k := range keys {
		command.append("-H", bashEscape(fmt.Sprintf("%s: %s", k, strings.Join(req.Header.GetAll(k), " "))))
	}

	command.append(bashEscape(req.URI().String()))

	return &command
}
