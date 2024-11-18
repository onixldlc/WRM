package display

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var (
	changeDisplaySettingsExW = user32.NewProc("ChangeDisplaySettingsExW")
)

const (
	CDS_UPDATEREGISTRY = 0x00000001
	CDS_TEST           = 0x00000002
)

// ChangeDisplaySettingsEx wraps the Windows API call
func ChangeDisplaySettingsEx(deviceName *uint16, lpDevMode *DEVMODE, hwnd uintptr, dwflags uint32, lParam uintptr) int32 {
	ret, _, _ := changeDisplaySettingsExW.Call(
		uintptr(unsafe.Pointer(deviceName)),
		uintptr(unsafe.Pointer(lpDevMode)),
		hwnd,
		uintptr(dwflags),
		lParam,
	)
	return int32(ret)
}

// SetResolution sets the resolution and frequency for a device
func SetResolution(deviceName string, resolution string, frequency uint32) error {
	modes, err := ListResolutions(deviceName)
	if err != nil {
		return err
	}
	resParts := strings.Split(resolution, "x")
	if len(resParts) != 2 {
		return fmt.Errorf("invalid resolution format. Use WidthxHeight (e.g., 1920x1080)")
	}
	width, err1 := strconv.Atoi(resParts[0])
	height, err2 := strconv.Atoi(resParts[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid resolution dimensions")
	}
	var selectedMode *DEVMODE
	for _, mode := range modes {
		if int(mode.DmPelsWidth) == width && int(mode.DmPelsHeight) == height {
			if frequency == 0 || mode.DmDisplayFrequency == frequency {
				if selectedMode == nil || mode.DmDisplayFrequency > selectedMode.DmDisplayFrequency {
					selectedMode = &mode
				}
			}
		}
	}
	if selectedMode == nil {
		return fmt.Errorf("resolution %s with frequency %d Hz not available", resolution, frequency)
	}
	// Confirm with the user
	fmt.Printf("Change resolution to %dx%d @ %d Hz? (y/n): ", width, height, selectedMode.DmDisplayFrequency)
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" {
		fmt.Println("Operation cancelled.")
		return nil
	}
	// Apply the settings with CDS_TEST flag first to validate
	deviceNamePtr, _ := syscall.UTF16PtrFromString(deviceName)
	result := ChangeDisplaySettingsEx(deviceNamePtr, selectedMode, 0, CDS_TEST, 0)
	if result != 0 {
		return fmt.Errorf("the requested graphics mode is not supported")
	}
	// Apply the settings and update the registry
	result = ChangeDisplaySettingsEx(deviceNamePtr, selectedMode, 0, CDS_UPDATEREGISTRY, 0)
	if result != 0 {
		return fmt.Errorf("failed to change display settings")
	}
	fmt.Println("Resolution changed successfully and saved to registry.")
	return nil
}
