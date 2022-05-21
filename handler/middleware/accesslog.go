package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ua "github.com/mileusna/useragent"
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
		Latency:   latency,
		Path:      path,
		OS:        os,
	}
}

func AccessLogger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		accessTimeBefore := time.Now()
		defer func() {
			ua := ua.Parse(r.UserAgent())
			accessTimeAfter := time.Now()
			accessTimeDiff := accessTimeAfter.Sub(accessTimeBefore).Microseconds()
			accessLog := NewAccessLog(accessTimeBefore, accessTimeDiff, r.URL.Path, ua.OS)
			accessLog.PrintJson()
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (a *AccessLog) PrintJson() {
	accessLogJson, err := json.Marshal(a)
	if err != nil {
		return
	}
	fmt.Println(string(accessLogJson))
}
