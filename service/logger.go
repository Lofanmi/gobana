package service

import (
	"context"
)

type Logger interface {
	// Search 日志搜索
	Search(ctx context.Context, req SearchRequest) (resp SearchResponse, err error)
}

type QueryBy interface {
	QueryByHuman | QueryByQueryString | QueryBySLS
}

type QueryByHuman struct {
	Or      []string `json:"or"`
	Must    []string `json:"must"`
	MustNot []string `json:"must_not"`
}

type QueryByQueryString struct {
	QueryString string `json:"query_string"`
}

type QueryBySLS struct {
	SQL    string `json:"sql"`
	Phrase string `json:"phrase"`
}

type SearchRequest struct {
	PageNo   int         `json:"page_no"`
	PageSize int         `json:"page_size"`
	TimeA    int64       `json:"time_a"`
	TimeB    int64       `json:"time_b"`
	Backend  string      `json:"backend"`
	Store    string      `json:"store"`
	QueryBy  int         `json:"query_by"`
	Query    interface{} `json:"query"`
}

type SearchResponse struct {
	PageNo   int         `json:"page_no"`
	PageSize int         `json:"page_size"`
	TimeA    int64       `json:"time_a"`
	TimeB    int64       `json:"time_b"`
	Count    int         `json:"count"`
	List     []LogItem   `json:"list"`
	RawQuery interface{} `json:"raw_query"`
}

type LogType = string

const (
	LogTypeAccessLog LogType = "access-log"
	LogTypeJsonLog   LogType = "json-log"
	LogTypeStringLog LogType = "string-log"
)

type AccessLog struct {
	Time          int64       `json:"time"`
	Method        string      `json:"method"`
	Scheme        string      `json:"scheme"`
	Host          string      `json:"host"`
	URI           string      `json:"uri"`
	Query         string      `json:"query"`
	Body          string      `json:"body"`
	Duration      string      `json:"duration"`
	HttpVersion   string      `json:"http_version"`
	UserAgent     string      `json:"user_agent"`
	Referer       string      `json:"referer"`
	XForwardedFor string      `json:"x_forwarded_for"`
	Cookie        string      `json:"cookie"`
	RemoteAddr    string      `json:"remote_addr"`
	CurlTemplate  string      `json:"curl_template"`
	Source        interface{} `json:"source"`
}

type JsonLog struct {
	Time     int64       `json:"time"`
	Level    string      `json:"level"`
	Hostname string      `json:"hostname"`
	Path     string      `json:"path"`
	Source   interface{} `json:"source"`
}

type StringLog struct {
	Time     int64       `json:"time"`
	Hostname string      `json:"hostname"`
	Path     string      `json:"path"`
	Source   interface{} `json:"source"`
}

type LogItem struct {
	Storage string      `json:"storage"`
	LogType LogType     `json:"log_type"`
	Log     interface{} `json:"log"`
}

type SearchCharts struct {
	Legend   []string           `json:"legend" form:"legend" query:"legend"`
	XAxis    []string           `json:"xAxis" form:"xAxis" query:"xAxis"`
	Series   SearchChartsSeries `json:"series" form:"series" query:"series"`
	Interval int                `json:"interval" form:"interval" query:"interval"`
}

type SearchChartsSeries struct {
	Name   string  `json:"name" form:"name" query:"name"`
	Type   string  `json:"type" form:"type" query:"type"`
	Symbol string  `json:"symbol" form:"symbol" query:"symbol"`
	Smooth bool    `json:"smooth" form:"smooth" query:"smooth"`
	Data   []int64 `json:"data" form:"data" query:"data"`
}

type SearchRequest2 struct {
	PageNo   int         `json:"page_no"`
	PageSize int         `json:"page_size"`
	TimeA    int64       `json:"time_a"`
	TimeB    int64       `json:"time_b"`
	Backend  string      `json:"backend"`
	Store    string      `json:"store"`
	QueryBy  int         `json:"query_by"`
	Query    interface{} `json:"query"`
}

type ExportRequest struct {
	Size    int64       `json:"size"`
	TimeA   int64       `json:"time_a"`
	TimeB   int64       `json:"time_b"`
	Backend string      `json:"backend"`
	Store   string      `json:"store"`
	QueryBy int         `json:"query_by"`
	Query   interface{} `json:"query"`
	Param   interface{} `json:"param"`
}

type ExportResponse struct {
	ID   string `json:"id"`
	Logs string `json:"logs"`
}
