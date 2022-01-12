package alice

type QueryRequest struct {
	Devices []struct {
		ID         string            `json:"id"`
		CustomData map[string]string `json:"custom_data"`
	} `json:"devices"`
}
