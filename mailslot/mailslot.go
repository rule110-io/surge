package mailslot

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	MAILSLOT_WAIT_FOREVER int32 = -1
	MAILSLOT_NO_MESSAGE   int32 = -1

	GENERIC_READ          = 0x80000000
	GENERIC_WRITE         = 0x40000000
	FILE_SHARE_READ       = 0x00000001
	FILE_SHARE_WRITE      = 0x00000002
	OPEN_EXISTING         = 3
	FILE_ATTRIBUTE_NORMAL = 0x00000080
)

//MailSlot .
type MailSlot struct {
	handle  uintptr
	timeout int32
}

//MailSlotFile .
type MailSlotFile struct {
	handle uintptr
}

//MailSlotInfo .
type MailSlotInfo struct {
	MaxSize     int32
	NextSize    int32
	Count       int32
	ReadTimeout int32
}

func (info MailSlotInfo) String() string {
	return fmt.Sprint("MaxMessageSize: ", info.MaxSize, ", NextSize: ", info.NextSize, ", Count: ", info.Count, ", ReadTimeout: ", info.ReadTimeout)
}

//New .
func New(name string, max int32, timeout int32) (*MailSlot, error) {
	handle, _, err := CreateMailslot(name, max, timeout)
	if handle == 0 {
		return nil, err
	}

	return &MailSlot{
		handle:  handle,
		timeout: timeout,
	}, nil
}

//Info .
func (ms *MailSlot) Info() (MailSlotInfo, error) {
	info := MailSlotInfo{}

	ok, _, err := GetMailslotInfo(ms.handle,
		uintptr(unsafe.Pointer(&info.MaxSize)),
		uintptr(unsafe.Pointer(&info.NextSize)),
		uintptr(unsafe.Pointer(&info.Count)),
		uintptr(unsafe.Pointer(&info.ReadTimeout)))

	if !ok {
		return info, err
	}
	return info, nil
}

//SetTimeout .
func (ms *MailSlot) SetTimeout(timeout int32) error {
	ok, _, err := SetMailslotInfo(ms.handle, timeout)
	if !ok {
		return err
	}
	ms.timeout = timeout
	return nil
}

//Read .
func (ms *MailSlot) Read(p []byte) (n int, err error) {
	return ReadFile(ms.handle, p)
}

//Close .
func (ms *MailSlot) Close() error {
	ok, _, err := CloseHandle(ms.handle)
	if ok {
		return nil
	}
	return err
}

//Open .
func Open(name string) (*MailSlotFile, error) {
	handle, err := CreateFile(name, GENERIC_READ|GENERIC_WRITE, FILE_SHARE_READ|FILE_SHARE_WRITE, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL)
	if int(handle) == -1 {
		return nil, err
	}
	return &MailSlotFile{handle: handle}, nil
}

func (ms *MailSlotFile) Write(p []byte) (n int, err error) {
	return WriteFile(ms.handle, p)
}

//Close .
func (ms *MailSlotFile) Close() error {
	ok, _, err := CloseHandle(ms.handle)
	if ok {
		return nil
	}
	return err
}

//SingleInstance checks if this client is the only instance
func SingleInstance(name string) error {
	ret, _, _ := createMailslot.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(`\\.\mailslot\`+name))),
		0,
		0,
		0,
	)
	// If the function fails, the return value is INVALID_HANDLE_VALUE.
	if int64(ret) == -1 {
		return fmt.Errorf("instance already exists")
	}
	return nil
}
