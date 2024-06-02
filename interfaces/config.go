package interfaces

import "time"

type Config interface {
	GetString(string) string
	Get(string) any
	GetInt(string) int
	GetDuration(string) time.Duration
	UnmarshalKey(s string, rawVal any) error
}
