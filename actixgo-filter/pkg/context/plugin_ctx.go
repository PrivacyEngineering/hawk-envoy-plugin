package context

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// NewHttpContext Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &HttpCtx{ContextID: contextID}
}
