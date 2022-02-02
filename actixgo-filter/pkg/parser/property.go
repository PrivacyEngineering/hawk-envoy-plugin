package parser

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

func getProperty(path []string) []byte {
	property, err := proxywasm.GetProperty(path)
	if err != nil {
		proxywasm.LogWarnf("Go-filter: property not available [%v]", path)
		return []byte("0")
	}
	return property
}
