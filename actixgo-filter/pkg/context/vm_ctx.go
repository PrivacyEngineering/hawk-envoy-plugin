package context

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type VmContext struct {
	types.DefaultVMContext
}

// NewPluginContext Override types.DefaultVMContext.
func (*VmContext) NewPluginContext(_ uint32) types.PluginContext {
	return &pluginContext{}
}
