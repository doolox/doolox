//go:build windows

package cms

import (
	"log"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")

	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procDestroyWindow    = user32.NewProc("DestroyWindow")
	procLoadIconW        = user32.NewProc("LoadIconW")
	procGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	procShellNotifyIconW = shell32.NewProc("Shell_NotifyIconW")
)

const (
	nimAdd    = 0x00000000
	nimDelete = 0x00000002

	nifIcon = 0x00000002
	nifInfo = 0x00000010

	niiInfo = 0x00000001

	idiApplication = 32512

	cwUseDefault = 0x80000000
	wsOverlapped = 0

	// HWND_MESSAGE (-3): a message-only window.
	hwndMessage = ^uintptr(2)
)

type wndClassExW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type notifyIconData struct {
	cbSize           uint32
	hWnd             uintptr
	uID              uint32
	uFlags           uint32
	uCallbackMessage uint32
	hIcon            uintptr
	szTip            [128]uint16
	dwState          uint32
	dwStateMask      uint32
	szInfo           [256]uint16
	uTimeout         uint32
	szInfoTitle      [64]uint16
	dwInfoFlags      uint32
	guidItem         [16]byte
	hBalloonIcon     uintptr
}

func utf16Ptr(s string) *uint16 {
	u := utf16.Encode([]rune(s))
	u = append(u, 0)
	return &u[0]
}

// notify shows a native Windows balloon notification via Shell_NotifyIcon.
// It is best-effort: failures are logged and otherwise ignored.
func notify(title, message string) {
	go func() {
		hInstance, _, _ := procGetModuleHandleW.Call(0)

		className := utf16Ptr("DooloxNotifyWindow")

		wc := wndClassExW{
			cbSize:        uint32(unsafe.Sizeof(wndClassExW{})),
			lpfnWndProc:   procDefWindowProcW.Addr(),
			hInstance:     hInstance,
			lpszClassName: className,
		}
		procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

		hWnd, _, _ := procCreateWindowExW.Call(
			0,
			uintptr(unsafe.Pointer(className)),
			uintptr(unsafe.Pointer(className)),
			wsOverlapped,
			cwUseDefault, cwUseDefault,
			cwUseDefault, cwUseDefault,
			hwndMessage,
			0,
			hInstance,
			0,
		)
		if hWnd == 0 {
			log.Printf("[notify] could not create notification window")
			return
		}
		defer procDestroyWindow.Call(hWnd)

		icon, _, _ := procLoadIconW.Call(0, uintptr(idiApplication))

		nid := notifyIconData{
			cbSize:      uint32(unsafe.Sizeof(notifyIconData{})),
			hWnd:        hWnd,
			uID:         1,
			uFlags:      nifIcon | nifInfo,
			hIcon:       icon,
			dwInfoFlags: niiInfo,
		}
		copy(nid.szInfoTitle[:], utf16.Encode([]rune(title)))
		copy(nid.szInfo[:], utf16.Encode([]rune(message)))

		procShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&nid)))

		// Keep the icon alive long enough for the balloon to display, then
		// remove it so nothing lingers in the system tray.
		time.Sleep(10 * time.Second)
		procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&nid)))
	}()
}
