//go:generate go run -tags generate gen.go
package main

import (
	"gosunxifel/config"
	"runtime"
	"syscall"

	log "github.com/donnie4w/go-logger/logger"
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	log.Debug("Running with CPUs:", nuCPU)
}

func AppVersion() {
	log.Debug("app-version:", config.Cfg.Section("").Key("app_ver").String())
}

func GetSystemMetrics(nIndex int) int {
	ret, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(nIndex))
	return int(ret)
}

func main() {
	ConfigRuntime()
	AppVersion()
	Defaultweb = New(GetSystemMetrics(0), GetSystemMetrics(1))
	Defaultweb.Run()
}
