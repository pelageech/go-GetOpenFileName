/*
	Most of the information took from
	https://learn.microsoft.com/en-us/windows/win32/api/commdlg/nf-commdlg-getopenfilenamew and
	https://learn.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-openfilenamew.

	The template took from https://go.dev/play/p/YjJydAov7G
	Thanks for such a helper.
*/

package ofn

import (
	"syscall"
	"unsafe"
)

const (
	OFN_ALLOWMULTISELECT     = 0x00000200
	OFN_CREATEPROMPT         = 0x00002000
	OFN_DONTADDTORECENT      = 0x02000000
	OFN_ENABLEHOOK           = 0x00000020
	OFN_ENABLEINCLUDENOTIFY  = 0x00400000
	OFN_ENABLESIZING         = 0x00800000
	OFN_ENABLETEMPLATE       = 0x00000040
	OFN_ENABLETEMPLATEHANDLE = 0x00000080
	OFN_EXPLORER             = 0x00080000
	OFN_EXTENSIONDIFFERENT   = 0x00000400
	OFN_FILEMUSTEXIST        = 0x00001000
	OFN_FORCESHOWHIDDEN      = 0x10000000
	OFN_HIDEREADONLY         = 0x00000004
	OFN_LONGNAMES            = 0x00200000
	OFN_NOCHANGEDIR          = 0x00000008
	OFN_NODEREFERENCELINKS   = 0x00100000
	OFN_NOLONGNAMES          = 0x00040000
	OFN_NONETWORKBUTTON      = 0x00020000
	OFN_NOREADONLYRETURN     = 0x00008000
	OFN_NOTESTFILECREATE     = 0x00010000
	OFN_NOVALIDATE           = 0x00000100
	OFN_OVERWRITEPROMPT      = 0x00000002
	OFN_PATHMUSTEXIST        = 0x00000800
	OFN_READONLY             = 0x00000001
	OFN_SHAREAWARE           = 0x00004000
	OFN_SHOWHELP             = 0x00000010
)

type WORD = uint16
type DWORD = uint32

var (
	dll                *syscall.DLL
	symGetOpenFileName *syscall.Proc
)

// Init loads comdlg32.dll into memory and looks for a symbol
// GetOpenFileNameW.
//
// Returns an error if dll loading fails of a symbol is not found.
func Init() error {
	var err error

	dll, err = syscall.LoadDLL("comdlg32.dll")
	if err != nil {
		return err
	}

	symGetOpenFileName, err = dll.FindProc("GetOpenFileNameW")
	if err != nil {
		return err
	}
	return nil
}

// OPENFILENAME describes parameters for GetOpenFileName.
// See https://learn.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-openfilenamew.
type OPENFILENAME struct {
	LStructSize       DWORD
	HwndOwner         syscall.Handle
	HInstance         syscall.Handle
	LpstrFilter       *uint16
	LpstrCustomFilter *uint16
	NMaxCustFilter    DWORD
	NFilterIndex      DWORD
	LpstrFile         *uint16
	NMaxFile          DWORD
	LpstrFileTitle    *uint16
	NMaxFileTitle     DWORD
	LpstrInitialDir   *uint16
	LpstrTitle        *uint16
	Flags             DWORD
	NFileOffset       WORD
	NFileExtension    WORD
	LpstrDefExt       *uint16
	LCustData         uintptr
	LpfnHook          uintptr
	LpTemplateName    *uint16
	PvReserved        uintptr
	DwReserved        DWORD
	FlagsEx           DWORD
}

// ChooseFileSimple is a simplified GetOpenFileName function
//
// It takes a string that configures LpstrFilter encoded to UTF-16 uint16 slice.
// It must approach to a format of Lpcstr (like "All Files\000*.*\000\000").
func ChooseFileSimple(lpcstr *uint16, flags DWORD, filePath []uint16) bool {
	var ofn OPENFILENAME

	ofn.LStructSize = DWORD(unsafe.Sizeof(ofn))

	//filePath := make([]uint16, 256)
	ofn.LpstrFile = &filePath[0]

	ofn.Flags = flags
	ofn.LpstrFilter = lpcstr
	ofn.NMaxFile = DWORD(len(filePath))

	return GetOpenFileName(&ofn)
}

// GetOpenFileName opens Windows Explorer where the user chooses
// one file to get a path to it.
// The function uses comdlg32.dll which contains GetOpenFileNameW symbol.
// It takes OPENFILENAME struct pointer.
//
// See more https://learn.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-openfilenamew.
func GetOpenFileName(lpofn *OPENFILENAME) bool {
	ret, _, _ := symGetOpenFileName.Call(uintptr(unsafe.Pointer(lpofn)))
	return ret != 0
}

// Release calls a release of comdlg32.dll.
// Returns an error if release is unsuccessful.
func Release() error {
	return dll.Release()
}
