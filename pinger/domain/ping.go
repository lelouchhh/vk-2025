package domain

import "time"

// PingResult represents the result of a ping operation.
type PingResult struct {
	IP          string    `json:"ip"`
	PingTime    float64   `json:"ping_time"`
	LastSuccess time.Time `json:"last_success"`
}

// Container represents a Docker container with its IP address.
type Container struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}
