package user32

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/christowolf/usb-event/internal/types"
)

// Taken from https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-registerdevicenotificationw#parameters.
const (
	DEVICE_NOTIFY_WINDOW_HANDLE = DEVICE_NOTIFY_HANDLE_TYPE(0x00000000)
	// Currently not supported: DEVICE_NOTIFY_SERVICE_HANDLE = DEVICE_NOTIFY_HANDLE_TYPE(0x00000001)
	DEVICE_NOTIFY_ALL_INTERFACE_CLASSES = types.DWORD(0x00000004)
	DBT_DEVTYP_DEVICEINTERFACE          = types.DWORD(0x00000005)
)

var (
	dll                             = syscall.NewLazyDLL("user32.dll")
	procRegisterDeviceNotificationW = dll.NewProc("RegisterDeviceNotificationW")
)

var (
	// classGuidUsb is the GUID for all USB serial host PnP drivers,
	// see https://learn.microsoft.com/en-us/windows/win32/devio/registering-for-device-notification?redirectedfrom=MSDN.
	classGuidUsb = windows.GUID{
		Data1: 0x25dbce51,
		Data2: 0x6c8f,
		Data3: 0x4a72,
		Data4: [8]byte{0x8a, 0x6d, 0xb5, 0x4c, 0x2b, 0x4f, 0xc8, 0x35},
	}
)

type DeviceNotificationFilter struct {
	size       types.DWORD
	deviceType types.DWORD
	reserved   types.DWORD
	classGuid  windows.GUID
	szName     uint16
}

type DEVICE_NOTIFY_HANDLE_TYPE types.DWORD

// RegisterDeviceNotificationW wraps
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-registerdevicenotificationw
// restricted to USB devices only.
func RegisterDeviceNotificationW(handle windows.Handle, flags DEVICE_NOTIFY_HANDLE_TYPE) error {
	var filter DeviceNotificationFilter
	filter.size = types.DWORD(unsafe.Sizeof(filter))
	filter.deviceType = DBT_DEVTYP_DEVICEINTERFACE
	filter.reserved = types.DWORD(0)
	filter.classGuid = classGuidUsb
	filter.szName = 0
	_, _, err := procRegisterDeviceNotificationW.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&filter)),
		uintptr(flags|DEVICE_NOTIFY_HANDLE_TYPE(DEVICE_NOTIFY_ALL_INTERFACE_CLASSES)))
	if !errors.Is(err, syscall.Errno(0)) {
		return err
	}
	return nil
}

// TODO: Is https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unregisterdevicenotification needed?
