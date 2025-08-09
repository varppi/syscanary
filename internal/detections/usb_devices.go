package detections

import (
	"fmt"
	"sync"
	"time"

	"github.com/Varppi/syscanary/internal/global"

	"github.com/google/gousb"
)

var (
	context     *gousb.Context
	deviceCount = -1
	lastDevices []string
)

func UsbDetect(output chan string, stop chan *struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	output <- "debug:usb:usb device monitoring started"
	context = gousb.NewContext()
	defer context.Close()
	for {
		outputData, err := check_state()
		if err != nil {
			output <- fmt.Sprintf("error:usb:%s", err.Error())
		}
		if outputData != "" {
			output <- fmt.Sprintf("info:usb:%s", outputData)
		}
		if len(stop) > 0 {
			break
		}
		time.Sleep(time.Duration(global.Config.Modules.Usb.Interval) * time.Second)
	}
}

func check_state() (string, error) {
	var devices []*gousb.DeviceDesc
	_, err := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		devices = append(devices, desc)
		return false
	})
	if err != nil {
		return "", err
	}

	if deviceCount == -1 {
		deviceCount = len(devices)
		findOutlier(devices) // Not actually trying to find any partciular device. It just updates the lastDevices variable :)
		return "", nil
	}

	relation := len(devices) - deviceCount
	if relation == 0 {
		return "", nil
	}

	status := "added"
	if relation < 0 {
		status = "removed"
	}

	outlier := findOutlier(devices)
	deviceCount = len(devices)
	return fmt.Sprintf("USB device %s (%s)", status, outlier), nil
}

func findOutlier(newDevices []*gousb.DeviceDesc) string {
	var tmpLastDevices []string
	var output = "<unknowwn device>"
	for _, newDevice := range newDevices {
		var found bool
		newDeviceIdentifier := newDevice.Product.String()
		tmpLastDevices = append(tmpLastDevices, newDeviceIdentifier)
		for _, lastDevice := range lastDevices {
			if newDeviceIdentifier == lastDevice {
				found = true
			}
		}

		// New device added
		if !found {
			output = newDeviceIdentifier
		}
	}

	for _, lastDevice := range lastDevices {
		var found bool
		lastDeviceIdentifier := lastDevice
		for _, newDevice := range newDevices {
			newDeviceIdentifier := newDevice.Product.String()
			if newDeviceIdentifier == lastDeviceIdentifier {
				found = true
			}
		}

		// Device removed
		if !found {
			output = lastDeviceIdentifier
		}
	}
	lastDevices = tmpLastDevices
	return output
}
