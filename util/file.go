package util

import (
	"os"
	"strings"

	log "github.com/donnie4w/go-logger/logger"
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDirectory(dirctory string) string {

	if dirctory == "" {
		return dirctory
	}
	if len(dirctory) == 4 && dirctory[2:] == "//" {
		return dirctory
	}
	if len(dirctory) == 1 && dirctory == "/" {
		return dirctory
	}
	if strings.LastIndex(dirctory, "/") == 0 {
		return "/"
	}

	if len(dirctory) > 4 { //like D://root

	}
	log.Debug(dirctory)
	re := dirctory[:strings.LastIndex(dirctory, `/`)]
	log.Debug(re)
	if len(re) == 3 {
		if string(dirctory[1]) == ":" {
			return re + "/"
		}
	}

	return re
}
