package model

import (
	"time"

	"github.com/Lofanmi/gobana/service"
)

type AccessLog struct {
	Time          time.Time
	Method        string
	Scheme        string
	Host          string
	URI           string
	Query         string
	Body          string
	Duration      time.Duration
	HttpVersion   string
	UserAgent     string
	Referer       string
	XForwardedFor string
	Cookie        string
	RemoteAddr    string
}

type JsonLog struct {
	Time     time.Time
	Level    string
	Hostname string
	Path     string
	Source   interface{}
}

type StringLog struct {
	Time     time.Time
	Hostname string
	Path     string
	Source   interface{}
}

type LogItem struct {
	Time    int64
	Storage string
	LogType service.LogType
	Log     interface{}
}

type LogItems []LogItem

func (s LogItems) Len() int           { return len(s) }
func (s LogItems) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s LogItems) Less(i, j int) bool { return s[i].Time > s[j].Time }
