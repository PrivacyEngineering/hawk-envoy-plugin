package context

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type VmContext struct {
	types.DefaultVMContext
}

// NewPluginContext Override types.DefaultVMContext.
func (*VmContext) NewPluginContext(vmCtxId uint32) types.PluginContext {
	configuration, err := proxywasm.GetVMConfiguration()
	if err != nil {
		proxywasm.LogWarnf("GO-filter: Not found vm configuration. Error:. %v", err)
		return &pluginContext{vmCtxId: vmCtxId}
	}
	vmConfig := string(configuration)
	proxywasm.LogWarnf("GO-filter: vm configuration found: %s", vmConfig)
	return &pluginContext{vmCtxId: vmCtxId, vmConfig: vmConfig}
}
