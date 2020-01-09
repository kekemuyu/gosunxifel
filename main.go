//go:generate go run -tags generate gen.go
package main

import (
	"gosunxifel/config"
	"runtime"

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

func main() {
	ConfigRuntime()
	AppVersion()
	Defaultweb = New(800, 600)
	Defaultweb.Run()
}
