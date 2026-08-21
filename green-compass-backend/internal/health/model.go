package health

type Status string

const StatusOK Status = "ok"

type StatusResponse struct {
	Status        Status `json:"status"`
	Version       string `json:"version"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}
