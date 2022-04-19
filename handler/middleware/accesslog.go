package middleware

import (
	"encoding/json"
	"fmt"
	"time"
)

type AccessLog struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency"`
	Path      string    `json:"path"`
	OS        string    `json:"os"`
}

func NewAccessLog(timestamp time.Time, latency int64, path, os string) *AccessLog {
	return &AccessLog{
		Timestamp: timestamp,
		Latency: latency,
		Path: path,
		OS: os,
	}
}

func (a *AccessLog) PrintAccessLogJson() {
	accessLogJson, err := json.Marshal(a)
	if err != nil {
		return
	}
	fmt.Println(string(accessLogJson))
}