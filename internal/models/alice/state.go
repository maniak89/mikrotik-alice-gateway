package alice

type State struct {
	TS      int64        `json:"ts"`
	Payload PayloadState `json:"payload"`
}

type PayloadState struct {
	UserID  string               `json:"user_id"`
	Devices []PayloadStateDevice `json:"devices"`
}

type PayloadStateDevice struct {
	ID         string                         `json:"id"`
	Properties []PayloadStateDeviceProperties `json:"properties"`
}

type PayloadStateDeviceProperties struct {
	Type  PropertyType                      `json:"type"`
	State PayloadStateDevicePropertiesState `json:"state"`
}

type PayloadStateDevicePropertiesState struct {
	Instance PropertyParameterInstance   `json:"instance"`
	Value    PropertyParameterValueValue `json:"value"`
}
