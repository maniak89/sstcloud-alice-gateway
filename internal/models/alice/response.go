package alice

type Response struct {
	RequestID string      `json:"request_id"`
	Payload   interface{} `json:"payload"`
}
