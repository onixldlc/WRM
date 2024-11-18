package display

import (
	"fmt"
	"strconv"
	"strings"
)

// ListFrequenciesForResolution lists available frequencies for a resolution
func ListFrequenciesForResolution(monitorIndex int, resolution string) {
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

	// Use the device name instead of the friendly name
	deviceName := mi.DeviceName

	modes, err := ListResolutions(deviceName)
	if err != nil {
		fmt.Println("Error listing frequencies:", err)
		return
	}
	resParts := strings.Split(resolution, "x")
	if len(resParts) != 2 {
		fmt.Println("Invalid resolution format. Use WidthxHeight (e.g., 1920x1080)")
		return
	}
	width, err1 := strconv.Atoi(resParts[0])
	height, err2 := strconv.Atoi(resParts[1])
	if err1 != nil || err2 != nil {
		fmt.Println("Invalid resolution dimensions")
		return
	}
	fmt.Printf("Frequencies for %dx%d on %s:\n", width, height, mi.FriendlyName)
	freqMap := make(map[uint32]bool)
	count := 1
	for _, mode := range modes {
		if int(mode.DmPelsWidth) == width && int(mode.DmPelsHeight) == height {
			freq := mode.DmDisplayFrequency
			if !freqMap[freq] {
				fmt.Printf("%d. %d Hz\n", count, freq)
				freqMap[freq] = true
				count++
			}
		}
	}
}

// ValidateFrequency checks if the frequency is valid for the given resolution on the monitor
func ValidateFrequency(deviceName string, resolution string, frequency uint32) (bool, error) {
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
		if int(mode.DmPelsWidth) == width && int(mode.DmPelsHeight) == height && mode.DmDisplayFrequency == frequency {
			return true, nil
		}
	}
	return false, nil
}

// GetHighestFrequency returns the highest frequency for the given resolution
func GetHighestFrequency(deviceName string, resolution string) (uint32, error) {
	modes, err := ListResolutions(deviceName)
	if err != nil {
		return 0, err
	}
	resParts := strings.Split(resolution, "x")
	if len(resParts) != 2 {
		return 0, fmt.Errorf("invalid resolution format")
	}
	width, err1 := strconv.Atoi(resParts[0])
	height, err2 := strconv.Atoi(resParts[1])
	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("invalid resolution dimensions")
	}
	var highestFreq uint32
	for _, mode := range modes {
		if int(mode.DmPelsWidth) == width && int(mode.DmPelsHeight) == height {
			if mode.DmDisplayFrequency > highestFreq {
				highestFreq = mode.DmDisplayFrequency
			}
		}
	}
	if highestFreq == 0 {
		return 0, fmt.Errorf("no frequencies found for resolution %s", resolution)
	}
	return highestFreq, nil
}
