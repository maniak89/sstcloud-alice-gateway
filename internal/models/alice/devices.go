package alice

type Devices struct {
	UserID  string   `json:"user_id"`
	Devices []Device `json:"devices"`
}
