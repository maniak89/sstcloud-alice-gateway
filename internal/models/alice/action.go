package alice

type ActionRequest struct {
	Payload struct {
		Devices []DeviceRequest `json:"devices"`
	} `json:"payload"`
}
