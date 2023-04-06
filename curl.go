package curl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cloudwego/hertz/cmd/hz/util"

	"github.com/cloudwego/hertz/pkg/protocol"
)

var ErrRequestURINotSet = fmt.Errorf("requestURI is not set")

// Command contains exec.Command compatible slice + helpers
type Command struct {
	slice []string
}

// Append appends a string to the CurlCommand
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

// GetCurlCommand returns a CurlCommand corresponding to an *protocol.Request
func GetCurlCommand(req *protocol.Request) (*Command, error) {
	fmt.Println(req.URI().String())
	if req.URI().String() == "http:///" {
		return nil, ErrRequestURINotSet
	}

	command := Command{}

	command.append("curl")

	if util.Bytes2Str(req.URI().Scheme()) == "https" {
		command.append("-k")
	}

	command.append("-X", bashEscape(util.Bytes2Str(req.Method())))

	body := req.Body()
	if len(body) > 0 {
		bodyEscaped := bashEscape(util.Bytes2Str(body))
		command.append("-d", bodyEscaped)
	}

	var keys []string

	req.Header.VisitAll(func(key, value []byte) {
		keys = append(keys, util.Bytes2Str(key))
	})
	sort.Strings(keys)

	for _, k := range keys {
		command.append("-H", bashEscape(fmt.Sprintf("%s: %s", k, strings.Join(req.Header.GetAll(k), " "))))
	}

	command.append(bashEscape(req.URI().String()))

	command.append("--compressed")

	return &command, nil
}
