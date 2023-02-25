package user32

import "syscall"

const (
	DEVICE_NOTIFY_WINDOW_HANDLE  = DEVICE_NOTIFY_HANDLE_TYPE(0x00000000)
	DEVICE_NOTIFY_SERVICE_HANDLE = DEVICE_NOTIFY_HANDLE_TYPE(0x00000001)
)

var (
	dll                             = syscall.NewLazyDLL("user32.dll")
	procRegisterDeviceNotificationW = dll.NewProc("RegisterDeviceNotificationW")
)

type DeviceNotificationFilter struct {
	size       DWORD
	deviceType DWORD
	reserved   DWORD
	classGuid  GUID
	szName     uint16
}

type DEVICE_NOTIFY_HANDLE_TYPE DWORD

type DWORD uint32

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type HANDLE uintptr

type LPVOID uintptr

func RegisterDeviceNotificationW(
	hRecipient HANDLE,
	notificationFilter LPVOID,
	flags DEVICE_NOTIFY_HANDLE_TYPE) HANDLE {
	ret, _, err := procRegisterDeviceNotificationW.Call(
		uintptr(hRecipient),
		uintptr(notificationFilter),
		uintptr(flags))
	if err != nil {
		panic(err)
	}
	return HANDLE(ret)
}
