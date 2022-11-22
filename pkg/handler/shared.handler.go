package handler

import "time"

const (
	refreshInterval = time.Duration(500) * time.Millisecond
)

type TurnSwitch struct {
	TurnOn bool `json:"turnOn"`
}
