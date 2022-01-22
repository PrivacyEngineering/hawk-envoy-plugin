package main

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// NewPluginContext Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(_ uint32) types.PluginContext {
	return &pluginContext{}
}
