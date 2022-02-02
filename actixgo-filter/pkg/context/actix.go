package context

import (
	"github.com/TUB-CNPE-TB/rust-envoy-proxy/actixgo-filter/pkg/parser"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

type Category int64

const (
	REQUEST  Category = 0
	RESPONSE          = 1
)

var (
	actixHeaders = [][2]string{
		{":method", "POST"},
		{":path", "/"},
		{":authority", "internal.org.net"},
		{"Content-Type", "application/json"},
	}
	//cluster  = "outbound|80||actix-collector-service.httpbin-gateway.svc.cluster.local"
	timeout  = 1000
	trailers [][2]string
	messages = make(map[uint32]Layer)
)

type Layer struct {
	RequestBody     []byte
	ResponseBody    []byte
	RequestHeaders  [][2]string
	ResponseHeaders [][2]string
}

func (ctx *HttpCtx) addLayer(body []byte, headers [][2]string, category Category) {
	layer := Layer{}
	if val, ok := messages[ctx.pluginCtxId]; ok {
		layer = val
	}

	if body != nil {
		if REQUEST == category {
			layer.RequestBody = body
		}
		if RESPONSE == category {
			layer.ResponseBody = body
		}
	}
	if headers != nil {
		if REQUEST == category {
			layer.RequestHeaders = headers
		}
		if RESPONSE == category {
			layer.ResponseHeaders = headers
		}
	}

	messages[ctx.pluginCtxId] = layer
}

func (ctx *HttpCtx) submit() {
	val, ok := messages[ctx.pluginCtxId]
	if !ok {
		proxywasm.LogWarnf("GO: No messages from context id %v", ctx.pluginCtxId)
		return
	}

	data := parser.Transform(val.RequestBody, val.ResponseBody, val.RequestHeaders, val.ResponseHeaders)
	ctx.callActix(data)
	delete(messages, ctx.pluginCtxId)
}

func (ctx *HttpCtx) callActix(data []byte) {
	proxywasm.LogWarnf("GO: filter in action to %s", ctx.vmConfig)
	_, err := proxywasm.DispatchHttpCall(ctx.vmConfig, actixHeaders, data, trailers, uint32(timeout), callback)
	if err != nil {
		proxywasm.LogWarnf("GO-callActix: Unable to perform action: %v", err)
	}
}

func callback(headers int, size int, trailers int) {
	proxywasm.LogWarnf("GO: Success with call. headers: %v, size: %v, trailers: %v", headers, size, trailers)
}
