package main

import (
	"github.com/TUB-CNPE-TB/rust-envoy-proxy/actixgo-filter/pkg/context"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

func main() {
	proxywasm.SetVMContext(&context.VmContext{})
}
