package capi

import "testing"

func TestGeneratedSymbolsExist(t *testing.T) {
	_ = Register
	_ = CEFExecuteProcess
	_ = CEFInitialize
	_ = CEFShutdown
	_ = CEFDoMessageLoopWork
	_ = CEFBrowserHostCreateBrowser
}
