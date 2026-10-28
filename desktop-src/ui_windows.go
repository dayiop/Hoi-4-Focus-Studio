//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

func showError(s string) {
	title, _ := syscall.UTF16PtrFromString("HOI4 Focus Studio")
	msg, _ := syscall.UTF16PtrFromString(s)
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	_, _, _ = proc.Call(0, uintptr(unsafe.Pointer(msg)), uintptr(unsafe.Pointer(title)), 0x10)
}
