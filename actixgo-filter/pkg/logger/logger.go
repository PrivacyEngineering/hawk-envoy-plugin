package logger

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"

func WriteLog(format string, args ...interface{}) {
	proxywasm.LogWarnf(format, args...)
}
