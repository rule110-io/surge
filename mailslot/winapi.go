package mailslot

import (
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	createMailslot  = kernel32.NewProc("CreateMailslotW")
	getMailslotInfo = kernel32.NewProc("GetMailslotInfo")
	setMailslotInfo = kernel32.NewProc("SetMailslotInfo")
	closeHandle     = kernel32.NewProc("CloseHandle")
	readFile        = kernel32.NewProc("ReadFile")
	createFile      = kernel32.NewProc("CreateFileW")
	writeFile       = kernel32.NewProc("WriteFile")
)

func CreateMailslot(name string, maxSize, readTimeout int32) (uintptr, int, error) {
	namePtr := uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name)))
	r1, r2, err := createMailslot.Call(namePtr, uintptr(maxSize), uintptr(readTimeout), 0)
	return r1, int(r2), err
}

func GetMailslotInfo(hMailslot uintptr, lpMaxMessageSize, lpNextSize, lpMessageCount, lpReadTimeout uintptr) (bool, int, error) {
	r1, r2, err := getMailslotInfo.Call(hMailslot, lpMaxMessageSize, lpNextSize, lpMessageCount, lpReadTimeout)
	return (r1 == uintptr(0)), int(r2), err
}

func SetMailslotInfo(hMailslot uintptr, lReadTimeout int32) (bool, int, error) {
	r1, r2, err := setMailslotInfo.Call(hMailslot, uintptr(lReadTimeout))
	return (r1 == uintptr(0)), int(r2), err
}

func CloseHandle(hObject uintptr) (bool, int, error) {
	r1, r2, err := closeHandle.Call(hObject)
	return (r1 == uintptr(0)), int(r2), err
}

func ReadFile(hFile uintptr, lpBuffer []byte) (int, error) {
	var rz int
	r1, _, err := readFile.Call(hFile,
		uintptr(unsafe.Pointer(&lpBuffer[0])),
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&rz)), 0)

	if r1 != 0 {
		return rz, nil
	}
	return 0, err
}

func CreateFile(name string, dwDesiredAccess, dwShareMode, dwCreationDisposition, dwFlagsAndAttributes int) (uintptr, error) {
	namePtr := uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name)))
	r1, _, err := createFile.Call(namePtr,
		uintptr(dwDesiredAccess),
		uintptr(dwShareMode),
		0,
		uintptr(dwCreationDisposition),
		uintptr(dwFlagsAndAttributes),
		0)
	if int(r1) == -1 {
		return r1, err
	}
	return r1, nil
}

func WriteFile(hFile uintptr, lpBuffer []byte) (int, error) {
	var sz int
	r1, _, err := writeFile.Call(hFile,
		uintptr(unsafe.Pointer(&lpBuffer[0])),
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&sz)),
		0)
	if int(r1) == 0 {
		return 0, err
	}
	return sz, nil
}
