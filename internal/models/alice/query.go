package alice

type QueryRequest struct {
	Devices []DeviceRequest `json:"devices"`
}
