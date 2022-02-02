package context

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type HttpCtx struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	vmCtxId     uint32
	pluginCtxId uint32
	vmConfig    string
}

// OnHttpRequestBody Override types.DefaultHttpContext.
func (ctx *HttpCtx) OnHttpRequestBody(bodySize int, _ bool) types.Action {
	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogWarnf("GO-OnHttpRequestBody: Unable to perform action: %v", err)
		return types.ActionContinue
	}
	ctx.addLayer(body, nil, REQUEST)
	return types.ActionContinue
}

// OnHttpResponseBody Override types.DefaultHttpContext.
func (ctx *HttpCtx) OnHttpResponseBody(bodySize int, _ bool) types.Action {
	body, err := proxywasm.GetHttpResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogWarnf("GO-OnHttpResponseBody: Unable to perform action: %v", err)
		return types.ActionContinue
	}
	ctx.addLayer(body, nil, RESPONSE)
	return types.ActionContinue
}

// OnHttpRequestHeaders Override types.DefaultHttpContext.
func (ctx *HttpCtx) OnHttpRequestHeaders(_ int, _ bool) types.Action {
	headers, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogWarnf("GO-OnHttpRequestHeaders: Unable to perform action: %v", err)
		return types.ActionContinue
	}
	ctx.addLayer(nil, headers, REQUEST)
	err = proxywasm.AddHttpRequestHeader("x-actix-go-filter", "processed")
	if err != nil {
		proxywasm.LogWarnf("Go-filter: Unable to add custome header x-actix-go-filter: %v", err)
	}
	return types.ActionContinue
}

// OnHttpResponseHeaders Override types.DefaultHttpContext.
func (ctx *HttpCtx) OnHttpResponseHeaders(_ int, _ bool) types.Action {
	headers, err := proxywasm.GetHttpResponseHeaders()
	if err != nil {
		proxywasm.LogWarnf("GO-OnHttpResponseHeaders: Unable to perform action: %v", err)
		return types.ActionContinue
	}
	ctx.addLayer(nil, headers, RESPONSE)
	return types.ActionContinue
}

// OnHttpStreamDone Override types.DefaultHttpContext.
func (ctx *HttpCtx) OnHttpStreamDone() {
	ctx.submit()
}
