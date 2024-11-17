package display

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type WmiMonitor struct {
	Name        string
	PNPDeviceID string
}

// GetWmiMonitors retrieves monitor names via WMI
func GetWmiMonitors() ([]WmiMonitor, error) {
	var monitors []WmiMonitor

	// Initialize COM
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// Connect to WMI
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, err
	}
	defer unknown.Release()

	wmiObj, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, err
	}
	defer wmiObj.Release()

	serviceRaw, err := oleutil.CallMethod(wmiObj, "ConnectServer")
	if err != nil {
		return nil, err
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// Execute WMI query
	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", "SELECT Name, PNPDeviceID FROM Win32_DesktopMonitor")
	if err != nil {
		return nil, err
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	// Iterate over the results
	countVar, err := oleutil.GetProperty(result, "Count")
	if err != nil {
		return nil, err
	}
	count := int(countVar.Val)

	for i := 0; i < count; i++ {
		itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
		if err != nil {
			return nil, err
		}
		item := itemRaw.ToIDispatch()
		defer item.Release()

		nameVar, err := oleutil.GetProperty(item, "Name")
		if err != nil {
			return nil, err
		}
		name := nameVar.ToString()

		pnpDeviceIDVar, err := oleutil.GetProperty(item, "PNPDeviceID")
		if err != nil {
			return nil, err
		}
		pnpDeviceID := pnpDeviceIDVar.ToString()

		monitors = append(monitors, WmiMonitor{
			Name:        name,
			PNPDeviceID: pnpDeviceID,
		})
	}

	return monitors, nil
}
