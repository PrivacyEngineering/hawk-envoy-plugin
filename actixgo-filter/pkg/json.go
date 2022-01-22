package main

import (
	"fmt"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

type Parser struct {
	data string
}

func (p *Parser) toJson(val Layer) []byte {
	p.data = `{"from": "go-filter", "content": {`
	p.parseHeaders("request_headers", val.RequestHeaders)
	p.parseHeaders("response_headers", val.ResponseHeaders)
	p.parseBody("request_body", val.RequestBody)
	p.parseBody("response_body", val.ResponseBody)
	p.data += `"__end__": true`
	p.data += "} }"

	proxywasm.LogWarnf("json: %s", p.data)

	return []byte(p.data)
}

func (p *Parser) parseHeaders(label string, headers [][2]string) {
	if headers != nil {
		size := len(headers)
		p.data += fmt.Sprintf(`"%s": {`, label)
		for i, header := range headers {
			p.data += fmt.Sprintf(`"%s": "%s"`, header[0], header[1])
			if i < size-1 {
				p.data += ","
			}
		}
		p.data += `},`
	}
}

func (p *Parser) parseBody(label string, body []byte) {
	if body != nil {
		p.data += fmt.Sprintf(`"%s": %s `, label, string(body))
		p.data += ","
	}
}
