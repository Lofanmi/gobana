package service

import (
	"context"
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
	SlsQuery string `json:"sls_query"`
}

type SearchRequest struct {
	PageNo         int    `json:"page_no"`
	PageSize       int    `json:"page_size"`
	TimeA          int64  `json:"time_a"`
	TimeB          int64  `json:"time_b"`
	Backend        string `json:"backend"`
	Storage        string `json:"storage"`
	QueryBy        string `json:"query_by"`
	Query          any    `json:"query"`
	ChartInterval  int    `json:"chart_interval"`
	ChartVisible   bool   `json:"chart_visible"`
	TrackTotalHits bool   `json:"track_total_hits"`
}

type SearchResponse struct {
	PageNo   int          `json:"page_no"`
	PageSize int          `json:"page_size"`
	TimeA    int64        `json:"time_a"`
	TimeB    int64        `json:"time_b"`
	Count    int          `json:"count"`
	List     LogItems     `json:"list"`
	Charts   SearchCharts `json:"charts"`
	RawQuery any          `json:"raw_query"`
}

type LogItem struct {
	Timestamp int64          `json:"-"`
	Storage   string         `json:"storage"`
	Source    map[string]any `json:"source"`
	Log       map[string]any `json:"log"`
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
	Size    int64  `json:"size"`
	TimeA   int64  `json:"time_a"`
	TimeB   int64  `json:"time_b"`
	Backend string `json:"backend"`
	Store   string `json:"store"`
	QueryBy int    `json:"query_by"`
	Query   any    `json:"query"`
	Param   any    `json:"param"`
}

type ExportResponse struct {
	ID   string `json:"id"`
	Logs string `json:"logs"`
}

type LogItems []LogItem

func (s LogItems) Len() int           { return len(s) }
func (s LogItems) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s LogItems) Less(i, j int) bool { return s[i].Timestamp > s[j].Timestamp }

// func curlTemplate(item *AccessLog) string {
// 	isJSONString := func(s string) bool {
// 		n := len(s)
// 		if n <= 1 {
// 			return false
// 		}
// 		s = strings.TrimSpace(s)
// 		if s[0] == '{' && s[n-1] == '}' {
// 			return gjson.Valid(s)
// 		}
// 		if s[0] == '[' && s[n-1] == ']' {
// 			return gjson.Valid(s)
// 		}
// 		return false
// 	}
// 	removeLineEnd := func(s string) string {
// 		return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\r", ""), "\n", " "))
// 	}
// 	s := fmt.Sprintf("curl -v -X '%s' -H 'Host: %s' \\\n", item.Method, item.HttpHost)
// 	switch item.Method {
// 	case http.MethodPost, http.MethodPut, http.MethodPatch:
// 		var contentType string
// 		if isJSONString(item.Body) {
// 			contentType = " -H 'Content-Type: application/json' \\\n"
// 		} else if strings.HasPrefix(item.Body, "<xml>") && strings.HasSuffix(item.Body, "</xml>") {
// 			contentType = " -H 'Content-Type: application/xml' \\\n"
// 		} else {
// 			contentType = " -H 'Content-Type: application/x-www-form-urlencoded' \\\n"
// 		}
// 		s += contentType
// 	}
// 	if item.UserAgent != "" {
// 		s += fmt.Sprintf(" -H 'User-Agent: %s' \\\n", removeLineEnd(item.UserAgent))
// 	}
// 	if item.Referer != "" {
// 		s += fmt.Sprintf(" -H 'Referer: %s' \\\n", removeLineEnd(item.Referer))
// 	}
// 	if item.Cookie != "" {
// 		s += fmt.Sprintf(" -H 'Cookie: %s' \\\n", removeLineEnd(item.Cookie))
// 	}
// 	if item.Body != "" {
// 		s += fmt.Sprintf(" -d '%s' \\\n", removeLineEnd(item.Body))
// 	}
// 	query := ""
// 	if item.Query != "" {
// 		query += "?" + item.Query
// 	}
// 	u, err := url.Parse(item.URI)
// 	if err != nil {
// 		return ""
// 	}
// 	return s + fmt.Sprintf(` '#SCHEME#://#HOST#%s%s'`, u.Path, query)
// }
