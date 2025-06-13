package domain

import "time"

// Event represents something that happened in the system
type Event interface {
	GetTimestamp() time.Time
	GetType() string
}
