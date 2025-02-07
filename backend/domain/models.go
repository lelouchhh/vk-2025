package domain

import "time"

type Account struct {
	ID       int    `db:"id" json:"id"`
	Login    string `db:"login" json:"login" validate:"required,min=3,max=255"`
	Password string `db:"password" json:"password" validate:"required,min=6,max=255"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type PingResult struct {
	IP          string    `json:"ip"`
	PingTime    float64   `json:"ping_time"`
	Status      bool      `json:"status"`
	LastSuccess time.Time `json:"last_success"`
}

type Container struct {
	ID        int       `db:"id" json:"id"`
	PingTime  float64   `db:"ping_time" json:"ping_time"`
	IPAddress string    `db:"ip_address" json:"ip_address"`
	LastPing  time.Time `db:"last_ping" json:"last_ping"`
	Status    bool      `db:"status" json:"status"`
}
