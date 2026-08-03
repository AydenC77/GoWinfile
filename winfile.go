package goWinFile

import (
	"errors"
	"fmt"
	"os"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modcomdlg32              = windows.NewLazySystemDLL("comdlg32.dll")
	procGetOpenFileNameW     = modcomdlg32.NewProc("GetOpenFileNameW")
	procCommDlgExtendedError = modcomdlg32.NewProc("CommDlgExtendedError")
)

const (
	ofnFileMustExist = 0x00001000
	ofnPathMustExist = 0x00000800
	ofnHideReadOnly  = 0x00000004
	ofnExplorer      = 0x00080000
	ofnNoChangeDir   = 0x00000008
)

type openFileNameW struct {
	lStructSize       uint32
	hwndOwner         uintptr
	hInstance         uintptr
	lpstrFilter       *uint16
	lpstrCustomFilter *uint16
	nMaxCustFilter    uint32
	nFilterIndex      uint32
	lpstrFile         *uint16
	nMaxFile          uint32
	lpstrFileTitle    *uint16
	nMaxFileTitle     uint32
	lpstrInitialDir   *uint16
	lpstrTitle        *uint16
	flags             uint32
	nFileOffset       uint16
	nFileExtension    uint16
	lpstrDefExt       *uint16
	lCustData         uintptr
	lpfnHook          uintptr
	lpTemplateName    *uint16
	pvReserved        unsafe.Pointer
	dwReserved        uint32
	flagsEx           uint32
}

func utf16NulPairs(pairs ...string) *uint16 {
	var all []uint16
	for _, s := range pairs {
		all = append(all, utf16.Encode([]rune(s))...)
		all = append(all, 0)
	}
	all = append(all, 0)
	return &all[0]
}

// usage:
// "Text Files (*.txt)", "*.txt"
func PickTextFile(filtername, filter string) (string, error) {
	if wd, err := os.Getwd(); err == nil {
		defer os.Chdir(wd)
	}

	fileBuf := make([]uint16, 32768)

	ofn := openFileNameW{
		lStructSize: uint32(unsafe.Sizeof(openFileNameW{})),
		lpstrFilter: utf16NulPairs(filtername, filter),
		lpstrFile:   &fileBuf[0],
		nMaxFile:    uint32(len(fileBuf)),
		flags:       ofnFileMustExist | ofnPathMustExist | ofnHideReadOnly | ofnExplorer | ofnNoChangeDir,
	}

	r, _, _ := procGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	if r == 0 {
		code, _, _ := procCommDlgExtendedError.Call()
		if code == 0 {
			return "", errors.New("dialog cancelled")
		}
		return "", fmt.Errorf("GetOpenFileNameW failed: 0x%x", code)
	}

	return windows.UTF16ToString(fileBuf), nil
}
