package mappers

import (
	"strings"

	"mikrotik-alice-gateway/internal/models/alice"
	"mikrotik-alice-gateway/internal/models/common"
)

func DeviceToAlice(device *common.Router) []alice.Device {
	result := make([]alice.Device, 0, len(device.Hosts))
	for _, host := range device.Hosts {
		result = append(result, DeviceHostToAlice(device, host))
	}
	return result
}

func DeviceHostToAliceEvent(device *common.Router, host *common.Host) alice.PayloadStateDevice {
	name := host.Name
	if name == "" {
		name = device.Name + "_" + host.ID
	}
	property := alice.PayloadStateDeviceProperties{
		Type: alice.PropertyTypeEvent,
		State: alice.PayloadStateDevicePropertiesState{
			Instance: alice.PropertyParameterInstanceGas,
			Value:    alice.PropertyParameterInstanceGasNotDetected,
		},
	}
	if host.Online {
		property.State.Value = alice.PropertyParameterInstanceGasDetected
	}
	return alice.PayloadStateDevice{
		ID: CreateHostDeviceID(device.ID, host.ID),
		Properties: []alice.PayloadStateDeviceProperties{
			property,
		},
	}

}

func DeviceHostToAlice(device *common.Router, host *common.Host) alice.Device {
	name := host.Name
	if name == "" {
		name = device.Name + "_" + host.ID
	}
	state := alice.PayloadStateDevicePropertiesState{
		Instance: alice.PropertyParameterInstanceGas,
		Value:    alice.PropertyParameterInstanceGasNotDetected,
	}
	if host.Online {
		state.Value = alice.PropertyParameterInstanceGasDetected
	}
	return alice.Device{
		ID:   CreateHostDeviceID(device.ID, host.ID),
		Name: name,
		Type: alice.DeviceTypeSensorGas,
		DeviceInfo: &alice.DeviceInfo{
			Model:        device.Model,
			SWVersion:    device.SWVersion,
			Manufacturer: device.Manufacturer,
		},
		CustomData: device.AdditionalFields,
		Properties: []alice.Property{
			{
				Type:        alice.PropertyTypeEvent,
				Retrievable: true,
				Reportable:  true,
				Parameters: []alice.PropertyParameter{
					{
						Instance: alice.PropertyParameterInstanceGas,
						Events: []alice.PropertyParameterValue{
							{
								Value: alice.PropertyParameterInstanceGasDetected,
							},
							{
								Value: alice.PropertyParameterInstanceGasNotDetected,
							},
						},
					},
				},
				State:          state,
				LastUpdated:    host.Updated,
				StateChangedAt: host.Changed,
			},
		},
	}
}

func ExtractDeviceID(str string) string {
	parts := strings.Split(str, "_")
	if len(parts) < 2 {
		return str
	}
	return parts[0]
}

func CreateHostDeviceID(str1, str2 string) string {
	return strings.Join([]string{str1, "host", str2}, "_")
}
