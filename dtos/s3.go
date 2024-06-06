package dtos

import "time"

type S3Object struct {
	Key          string
	LastModified time.Time
	Size         int64
}
