package context

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	vmCtxId  uint32
	vmConfig string
	types.DefaultPluginContext
}

// NewHttpContext Override types.DefaultPluginContext.
func (p *pluginContext) NewHttpContext(pluginCtxId uint32) types.HttpContext {
	proxywasm.LogWarnf("Go-filter: creating new plugin context: %d", pluginCtxId)
	return &HttpCtx{vmConfig: p.vmConfig, vmCtxId: p.vmCtxId, pluginCtxId: pluginCtxId}
}
