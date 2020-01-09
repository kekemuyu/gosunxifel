package config

import (
	"os"
	"path/filepath"
	"strings"

	"gofile/util"

	log "github.com/donnie4w/go-logger/logger"
	"github.com/go-ini/ini"
)

var Cfg *ini.File

func init() {
	var err error

	Cfg, err = ini.Load(GetRootdir() + "/config/conf.ini")

	if err != nil {
		log.Error("Fail to read file: %v", err)
		os.Exit(1)
	}
}

func GetRootdir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		if util.Exist(d + "/config") {
			return d
		}

		return infer(filepath.Dir(d))
	}
	var re string
	re = strings.Replace(infer(cwd), `\`, `/`, -1)
	re = strings.Replace(re, `:/`, `://`, -1)
	log.Debug(re)
	return re
}

func Save() {
	Cfg.SaveTo(GetRootdir() + "/config/conf.ini")
}
