package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

type Logger interface {
	// Search 日志搜索
	Search(ctx context.Context, req SearchRequest) (resp SearchResponse, err error)
}

type QueryType = string

const (
	QueryTypeByHuman    QueryType = "query_by_human"
	QueryTypeByLucene   QueryType = "query_by_lucene"
	QueryTypeBySLSQuery QueryType = "query_by_sls_query"
)

type QueryByHuman struct {
	Or      []string `json:"or"`
	Must    []string `json:"must"`
	MustNot []string `json:"must_not"`
}

type QueryByLucene struct {
	Lucene string `json:"lucene"`
}

type QueryBySLSQuery struct {
	SQL string `json:"sql"`
}

type SearchRequest struct {
	PageNo         int         `json:"page_no"`
	PageSize       int         `json:"page_size"`
	TimeA          int64       `json:"time_a"`
	TimeB          int64       `json:"time_b"`
	Backend        string      `json:"backend"`
	Storage        string      `json:"storage"`
	QueryBy        string      `json:"query_by"`
	Query          interface{} `json:"query"`
	ChartInterval  int32       `json:"chart_interval"`
	ChartVisible   bool        `json:"chart_visible"`
	TrackTotalHits bool        `json:"track_total_hits"`
}

type SearchResponse struct {
	PageNo   int          `json:"page_no"`
	PageSize int          `json:"page_size"`
	TimeA    int64        `json:"time_a"`
	TimeB    int64        `json:"time_b"`
	Count    int          `json:"count"`
	List     []LogItem    `json:"list"`
	Charts   SearchCharts `json:"charts"`
	RawQuery interface{}  `json:"raw_query"`
}

type LogType = string

const (
	LogTypeAccessLog LogType = "access-log"
	LogTypeJsonLog   LogType = "json-log"
	LogTypeStringLog LogType = "string-log"
)

type Log interface {
	*AccessLog | *JsonLog | *StringLog
	GetSource() interface{}
	SetSource(v interface{})
	Finish()
}

type AccessLog struct {
	RequestID     string      `json:"request_id"`
	Time          string      `json:"time"`
	Method        string      `json:"method"`
	Scheme        string      `json:"scheme"`
	Hostname      string      `json:"hostname"`
	URI           string      `json:"uri"`
	HttpHost      string      `json:"http_host"`
	Query         string      `json:"query"`
	Body          string      `json:"body"`
	Duration      string      `json:"duration"`
	HttpVersion   string      `json:"http_version"`
	UserAgent     string      `json:"user_agent"`
	Referer       string      `json:"referer"`
	XForwardedFor string      `json:"x_forwarded_for"`
	Cookie        string      `json:"cookie"`
	RemoteAddr    string      `json:"remote_addr"`
	Status        int         `json:"status"`
	Message       string      `json:"message"`
	CurlTemplate  string      `json:"curl_template"`
	Source        interface{} `json:"source"`
}

type JsonLog struct {
	RequestID string      `json:"request_id"`
	Time      string      `json:"time"`
	Level     string      `json:"level"`
	Hostname  string      `json:"hostname"`
	Path      string      `json:"path"`
	Tag       string      `json:"tag"`
	Message   string      `json:"message"`
	Source    interface{} `json:"source"`
}

type StringLog struct {
	Time     string      `json:"time"`
	Hostname string      `json:"hostname"`
	Path     string      `json:"path"`
	Message  string      `json:"message"`
	Source   interface{} `json:"source"`
}

type LogItem struct {
	Timestamp int64       `json:"-"`
	Storage   string      `json:"storage"`
	LogType   LogType     `json:"log_type"`
	Log       interface{} `json:"log"`
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

func (s *AccessLog) GetSource() interface{}  { return s.Source }
func (s *AccessLog) SetSource(v interface{}) { s.Source = v }
func (s *JsonLog) GetSource() interface{}    { return s.Source }
func (s *JsonLog) SetSource(v interface{})   { s.Source = v }
func (s *StringLog) GetSource() interface{}  { return s.Source }
func (s *StringLog) SetSource(v interface{}) { s.Source = v }

func (s *AccessLog) Finish() {
	s.Time = formatTime(s.Time)
	if s.Scheme == "" {
		s.Scheme = "http"
	}
	if !strings.HasPrefix(s.URI, "http") {
		s.URI = strings.Trim(s.URI, ":/")
		s.URI = s.Scheme + "://" + s.URI
	}
	u, err := url.Parse(s.URI)
	if err != nil {
		return
	}
	s.HttpHost = u.Hostname()
	s.Query = u.RawQuery
	s.Duration = formatDuration(s.Duration)
	s.CurlTemplate = curlTemplate(s)
}

func (s *JsonLog) Finish() {
	s.Time = formatTime(s.Time)
	s.Level = strings.ToLower(s.Level)
}

func (s *StringLog) Finish() {
	s.Time = formatTime(s.Time)
}

type LogItems []LogItem

func (s LogItems) Len() int           { return len(s) }
func (s LogItems) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s LogItems) Less(i, j int) bool { return s[i].Timestamp > s[j].Timestamp }

func formatTime(s string) (res string) {
	s = strings.TrimSpace(s)
	if regexp.MustCompile(`^\d+$`).MatchString(s) {
		i, _ := strconv.Atoi(s)
		if len(s) == 10 {
			i *= 1000
		}
		return time.UnixMilli(int64(i)).Format(time.RFC3339Nano)
	}
	if strings.HasSuffix(s, "Z") || strings.Contains(s, "+") {
		t, err := time.ParseInLocation(time.RFC3339, s, time.Local)
		if err == nil {
			return t.Format(time.RFC3339Nano)
		}
	}
	t, err := time.Parse("2006-01-02 15:04:05.000000Z07:00", s)
	if err == nil {
		return t.Format(time.RFC3339Nano)
	}
	t, err = time.Parse("2006-01-02 15:04:05.000Z07:00", s)
	if err == nil {
		return t.Format(time.RFC3339Nano)
	}
	t, err = cast.ToTimeInDefaultLocationE(s, time.Local)
	if err == nil {
		return t.Format(time.RFC3339Nano)
	}
	return s
}

func formatDuration(s string) (res string) {
	duration, err := strconv.ParseFloat(s, 64)
	if err != nil {
		res = s
		return
	}
	t := time.Duration(duration * float64(time.Second))
	res = t.String()
	return
}

func curlTemplate(item *AccessLog) string {
	isJSONString := func(s string) bool {
		n := len(s)
		if n <= 1 {
			return false
		}
		s = strings.TrimSpace(s)
		if s[0] == '{' && s[n-1] == '}' {
			return gjson.Valid(s)
		}
		if s[0] == '[' && s[n-1] == ']' {
			return gjson.Valid(s)
		}
		return false
	}
	removeLineEnd := func(s string) string {
		return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\r", ""), "\n", " "))
	}
	s := fmt.Sprintf("curl -v -X '%s' -H 'Host: %s' \\\n", item.Method, item.HttpHost)
	switch item.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		var contentType string
		if isJSONString(item.Body) {
			contentType = " -H 'Content-Type: application/json' \\\n"
		} else if strings.HasPrefix(item.Body, "<xml>") && strings.HasSuffix(item.Body, "</xml>") {
			contentType = " -H 'Content-Type: application/xml' \\\n"
		} else {
			contentType = " -H 'Content-Type: application/x-www-form-urlencoded' \\\n"
		}
		s += contentType
	}
	if item.UserAgent != "" {
		s += fmt.Sprintf(" -H 'User-Agent: %s' \\\n", removeLineEnd(item.UserAgent))
	}
	if item.Referer != "" {
		s += fmt.Sprintf(" -H 'Referer: %s' \\\n", removeLineEnd(item.Referer))
	}
	if item.Cookie != "" {
		s += fmt.Sprintf(" -H 'Cookie: %s' \\\n", removeLineEnd(item.Cookie))
	}
	if item.Body != "" {
		s += fmt.Sprintf(" -d '%s' \\\n", removeLineEnd(item.Body))
	}
	query := ""
	if item.Query != "" {
		query += "?" + item.Query
	}
	u, err := url.Parse(item.URI)
	if err != nil {
		return ""
	}
	return s + fmt.Sprintf(` '#SCHEME#://#HOST#%s%s'`, u.Path, query)
}
