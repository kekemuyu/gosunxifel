// +build windows

package util

import (
	"syscall"
	"unsafe"
)

//硬盘信息
func GetDiskInfo() (infos []string) {
	kernel := syscall.NewLazyDLL("Kernel32.dll")
	GetLogicalDriveStringsW := kernel.NewProc("GetLogicalDriveStringsW")

	lpBuffer := make([]byte, 254)
	diskret, _, _ := GetLogicalDriveStringsW.Call(
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&lpBuffer[0])))
	if diskret == 0 {
		return
	}

	for _, v := range lpBuffer {
		if v >= 65 && v <= 90 {
			path := string(v) + "://"
			if path == "A:" || path == "B:" {
				continue
			}
			infos = append(infos, path)
		}
	}
	return infos
}
