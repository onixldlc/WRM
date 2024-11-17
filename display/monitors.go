package display

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32                          = syscall.NewLazyDLL("user32.dll")
	procGetDisplayConfigBufferSizes = user32.NewProc("GetDisplayConfigBufferSizes")
	procQueryDisplayConfig          = user32.NewProc("QueryDisplayConfig")
	procDisplayConfigGetDeviceInfo  = user32.NewProc("DisplayConfigGetDeviceInfo")
)

const (
	QDC_ONLY_ACTIVE_PATHS                     = 0x00000002
	DISPLAYCONFIG_DEVICE_INFO_GET_TARGET_NAME = 0x00000002
	DISPLAYCONFIG_DEVICE_INFO_GET_SOURCE_NAME = 0x00000001
	ERROR_SUCCESS                             = 0
)

type LUID struct {
	LowPart  uint32
	HighPart int32
}

type DISPLAYCONFIG_PATH_INFO struct {
	SourceInfo DISPLAYCONFIG_PATH_SOURCE_INFO
	TargetInfo DISPLAYCONFIG_PATH_TARGET_INFO
	Flags      uint32
}

type DISPLAYCONFIG_PATH_SOURCE_INFO struct {
	AdapterId   LUID
	Id          uint32
	ModeInfoIdx uint32
	StatusFlags uint32
}

type DISPLAYCONFIG_PATH_TARGET_INFO struct {
	AdapterId        LUID
	Id               uint32
	ModeInfoIdx      uint32
	OutputTechnology uint32
	Rotation         uint32
	Scaling          uint32
	RefreshRate      DISPLAYCONFIG_RATIONAL
	ScanLineOrdering uint32
	TargetAvailable  uint32
	StatusFlags      uint32
}

type DISPLAYCONFIG_RATIONAL struct {
	Numerator   uint32
	Denominator uint32
}

type DISPLAYCONFIG_MODE_INFO struct {
	InfoType  uint32
	Id        uint32
	AdapterId LUID
	// The Union field is represented as a byte array
	Union [64]byte
}

type DISPLAYCONFIG_DEVICE_INFO_HEADER struct {
	Type      uint32
	Size      uint32
	AdapterId LUID
	Id        uint32
}

type DISPLAYCONFIG_TARGET_DEVICE_NAME_FLAGS struct {
	Value uint32
}

type DISPLAYCONFIG_TARGET_DEVICE_NAME struct {
	Header                    DISPLAYCONFIG_DEVICE_INFO_HEADER
	Flags                     DISPLAYCONFIG_TARGET_DEVICE_NAME_FLAGS
	OutputTechnology          uint32
	EdidManufactureId         uint16
	EdidProductCodeId         uint16
	ConnectorInstance         uint32
	MonitorFriendlyDeviceName [64]uint16
	MonitorDevicePath         [128]uint16
}

type DISPLAYCONFIG_SOURCE_DEVICE_NAME struct {
	Header            DISPLAYCONFIG_DEVICE_INFO_HEADER
	ViewGdiDeviceName [32]uint16
}

func GetDisplayConfigBufferSizes(flags uint32, numPathArrayElements *uint32, numModeInfoArrayElements *uint32) int32 {
	ret, _, _ := procGetDisplayConfigBufferSizes.Call(
		uintptr(flags),
		uintptr(unsafe.Pointer(numPathArrayElements)),
		uintptr(unsafe.Pointer(numModeInfoArrayElements)),
	)
	return int32(ret)
}

func QueryDisplayConfig(flags uint32, numPathArrayElements *uint32, pathArray *DISPLAYCONFIG_PATH_INFO, numModeInfoArrayElements *uint32, modeInfoArray *DISPLAYCONFIG_MODE_INFO, currentTopologyId *uint32) int32 {
	ret, _, _ := procQueryDisplayConfig.Call(
		uintptr(flags),
		uintptr(unsafe.Pointer(numPathArrayElements)),
		uintptr(unsafe.Pointer(pathArray)),
		uintptr(unsafe.Pointer(numModeInfoArrayElements)),
		uintptr(unsafe.Pointer(modeInfoArray)),
		uintptr(unsafe.Pointer(currentTopologyId)),
	)
	return int32(ret)
}

func DisplayConfigGetDeviceInfo(requestPacket *DISPLAYCONFIG_DEVICE_INFO_HEADER) int32 {
	ret, _, _ := procDisplayConfigGetDeviceInfo.Call(
		uintptr(unsafe.Pointer(requestPacket)),
	)
	return int32(ret)
}

type MonitorInfo struct {
	AdapterId    LUID
	Id           uint32
	FriendlyName string
	DeviceName   string // e.g., "\\.\DISPLAY1"
}

// GetSourceDeviceName retrieves the source device name for the monitor
func GetSourceDeviceName(adapterId LUID, id uint32) (string, error) {
	var deviceName DISPLAYCONFIG_SOURCE_DEVICE_NAME
	deviceName.Header.Type = DISPLAYCONFIG_DEVICE_INFO_GET_SOURCE_NAME
	deviceName.Header.Size = uint32(unsafe.Sizeof(deviceName))
	deviceName.Header.AdapterId = adapterId
	deviceName.Header.Id = id

	ret := DisplayConfigGetDeviceInfo(&deviceName.Header)
	if ret != ERROR_SUCCESS {
		return "", fmt.Errorf("DisplayConfigGetDeviceInfo failed with error %d", ret)
	}

	return syscall.UTF16ToString(deviceName.ViewGdiDeviceName[:]), nil
}

// ListMonitors retrieves all active monitors with their friendly names and device names
func ListMonitors() ([]MonitorInfo, error) {
	var pathCount, modeCount uint32

	// Get buffer sizes
	ret := GetDisplayConfigBufferSizes(QDC_ONLY_ACTIVE_PATHS, &pathCount, &modeCount)
	if ret != ERROR_SUCCESS {
		return nil, fmt.Errorf("GetDisplayConfigBufferSizes failed with error %d", ret)
	}

	// Allocate the path and mode arrays
	pathArray := make([]DISPLAYCONFIG_PATH_INFO, pathCount)
	modeInfoArray := make([]DISPLAYCONFIG_MODE_INFO, modeCount)

	// Query display config
	ret = QueryDisplayConfig(QDC_ONLY_ACTIVE_PATHS, &pathCount, &pathArray[0], &modeCount, &modeInfoArray[0], nil)
	if ret != ERROR_SUCCESS {
		return nil, fmt.Errorf("QueryDisplayConfig failed with error %d", ret)
	}

	var monitors []MonitorInfo

	// Iterate over the paths
	for i := 0; i < int(pathCount); i++ {
		path := pathArray[i]

		// Prepare the DISPLAYCONFIG_TARGET_DEVICE_NAME structure
		var targetName DISPLAYCONFIG_TARGET_DEVICE_NAME
		targetName.Header.Type = DISPLAYCONFIG_DEVICE_INFO_GET_TARGET_NAME
		targetName.Header.Size = uint32(unsafe.Sizeof(targetName))
		targetName.Header.AdapterId = path.TargetInfo.AdapterId
		targetName.Header.Id = path.TargetInfo.Id

		// Get the device info
		ret = DisplayConfigGetDeviceInfo(&targetName.Header)
		if ret != ERROR_SUCCESS {
			fmt.Printf("DisplayConfigGetDeviceInfo failed with error %d\n", ret)
			continue
		}

		// Get the monitor friendly name
		friendlyName := syscall.UTF16ToString(targetName.MonitorFriendlyDeviceName[:])
		if friendlyName == "" {
			friendlyName = "Unknown Monitor"
		}

		// Get the source device name
		sourceDeviceName, err := GetSourceDeviceName(path.SourceInfo.AdapterId, path.SourceInfo.Id)
		if err != nil {
			fmt.Printf("GetSourceDeviceName failed with error: %v\n", err)
			continue
		}

		// The device name might be like "DISPLAY1", we need to prepend "\\.\"
		fullDeviceName := sourceDeviceName

		monitors = append(monitors, MonitorInfo{
			AdapterId:    path.SourceInfo.AdapterId,
			Id:           path.SourceInfo.Id,
			FriendlyName: friendlyName,
			DeviceName:   fullDeviceName,
		})
	}

	return monitors, nil
}

// PrintMonitors lists all monitors with their friendly names and device names
func PrintMonitors() error {
	monitors, err := ListMonitors()
	if err != nil {
		return err
	}
	for i, mi := range monitors {
		fmt.Printf("%d. %s (%s)\n", i+1, mi.FriendlyName, mi.DeviceName)
	}
	return nil
}
