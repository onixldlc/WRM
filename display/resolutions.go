package display

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// DEVMODE structure
type DEVMODE struct {
	DmDeviceName         [32]uint16
	DmSpecVersion        uint16
	DmDriverVersion      uint16
	DmSize               uint16
	DmDriverExtra        uint16
	DmFields             uint32
	DmPosition           POINTL
	DmDisplayOrientation uint32
	DmDisplayFixedOutput uint32
	DmColor              uint16
	DmDuplex             uint16
	DmYResolution        uint16
	DmTTOption           uint16
	DmCollate            uint16
	DmFormName           [32]uint16
	DmLogPixels          uint16
	DmBitsPerPel         uint32
	DmPelsWidth          uint32
	DmPelsHeight         uint32
	DmDisplayFlags       uint32
	DmDisplayFrequency   uint32
	DmICMMethod          uint32
	DmICMIntent          uint32
	DmMediaType          uint32
	DmDitherType         uint32
	DmReserved1          uint32
	DmReserved2          uint32
	DmPanningWidth       uint32
	DmPanningHeight      uint32
}

// POINTL structure
type POINTL struct {
	X int32
	Y int32
}

var (
	enumDisplaySettingsExW = user32.NewProc("EnumDisplaySettingsExW")
)

// Resolution represents a display resolution
type Resolution struct {
	Width  uint32
	Height uint32
}

// ListResolutions lists all available modes for a device
func ListResolutions(deviceName string) ([]DEVMODE, error) {
	var modes []DEVMODE
	var iModeNum uint32 = 0
	deviceNamePtr, _ := syscall.UTF16PtrFromString(deviceName)
	for {
		var devMode DEVMODE
		devMode.DmSize = uint16(unsafe.Sizeof(devMode))
		ret, _, _ := enumDisplaySettingsExW.Call(
			uintptr(unsafe.Pointer(deviceNamePtr)),
			uintptr(iModeNum),
			uintptr(unsafe.Pointer(&devMode)),
			0,
		)
		if ret == 0 {
			break
		}
		modes = append(modes, devMode)
		iModeNum++
	}
	return modes, nil
}

// ListResolutionsForMonitor lists available resolutions for a monitor
func ListResolutionsForMonitor(monitorIndex int) {
	monitors, err := ListMonitors()
	if err != nil {
		fmt.Println("Error listing monitors:", err)
		return
	}
	if monitorIndex < 0 || monitorIndex >= len(monitors) {
		fmt.Println("Monitor index out of range")
		return
	}
	mi := monitors[monitorIndex]

	// Use the device name directly
	deviceName := mi.DeviceName

	modes, err := ListResolutions(deviceName)
	if err != nil {
		fmt.Println("Error listing resolutions:", err)
		return
	}
	fmt.Printf("Resolutions for %s (%s):\n", mi.FriendlyName, mi.DeviceName)

	// Collect unique resolutions
	resolutionMap := make(map[string]Resolution)
	for _, mode := range modes {
		resKey := fmt.Sprintf("%dx%d", mode.DmPelsWidth, mode.DmPelsHeight)
		resolutionMap[resKey] = Resolution{
			Width:  mode.DmPelsWidth,
			Height: mode.DmPelsHeight,
		}
	}

	// Create a slice to sort resolutions
	var resolutions []Resolution
	for _, res := range resolutionMap {
		resolutions = append(resolutions, res)
	}

	// // Sort resolutions from highest to lowest (width * height)
	// sort.Slice(resolutions, func(i, j int) bool {
	// 	pixelsI := resolutions[i].Width * resolutions[i].Height
	// 	pixelsJ := resolutions[j].Width * resolutions[j].Height
	// 	if pixelsI == pixelsJ {
	// 		return resolutions[i].Width > resolutions[j].Width
	// 	}
	// 	return pixelsI > pixelsJ
	// })

	// Sort resolutions by their string representation in descending order
	sort.Slice(resolutions, func(i, j int) bool {
		// Create string representations
		strI := fmt.Sprintf("%05dx%05d", resolutions[i].Width, resolutions[i].Height)
		strJ := fmt.Sprintf("%05dx%05d", resolutions[j].Width, resolutions[j].Height)

		// Compare strings
		if strI == strJ {
			return false // They are equal; maintain original order
		}
		return strI > strJ // For descending order
	})

	// Print sorted resolutions
	for i, res := range resolutions {
		fmt.Printf("%3d. %dx%d\n", i+1, res.Width, res.Height)
	}
}

// ValidateResolution checks if the resolution is valid for the given monitor
func ValidateResolution(deviceName string, resolution string) (bool, error) {
	modes, err := ListResolutions(deviceName)
	if err != nil {
		return false, err
	}
	resParts := strings.Split(resolution, "x")
	if len(resParts) != 2 {
		return false, fmt.Errorf("invalid resolution format")
	}
	width, err1 := strconv.Atoi(resParts[0])
	height, err2 := strconv.Atoi(resParts[1])
	if err1 != nil || err2 != nil {
		return false, fmt.Errorf("invalid resolution dimensions")
	}
	for _, mode := range modes {
		if int(mode.DmPelsWidth) == width && int(mode.DmPelsHeight) == height {
			return true, nil
		}
	}
	return false, nil
}
