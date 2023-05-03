package svc_logger

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/service"
)

const (
	atTimestamp    = "@timestamp"
	lenChartPoints = 60 // 图表固定最多60个数据点
)

func (s Service) SearchSLSLog(ctx context.Context, req service.LoggerSearchRequest) (items service.LoggerLogItems, count int, xAxis []string, yAxis []int64, err error) {
	return s.searchSlsLog(ctx, req.ChartsVisible, req.TimeA/1e3, req.TimeB/1e3, req.PageSize, req.PageNo, req.Index, req.Phrase, req.Sql)
}

func (s Service) BuildSlsQuery(ctx context.Context, or, must, mustNot []string) (query string) {
	mustNot = append(mustNot, "sy-job-manager.39on.com")
	orQuery := strings.Join(or, " OR ")
	mustQuery := strings.Join(must, " AND ")
	mustNotQuery := strings.Join(mustNot, " AND ")
	q := make([]string, 0)
	if len(orQuery) > 0 {
		q = append(q, "("+orQuery+")")
	}
	if len(mustQuery) > 0 {
		q = append(q, "("+mustQuery+")")
	}
	if len(mustNotQuery) > 0 {
		q = append(q, "NOT ("+mustNotQuery+")")
	}
	query = strings.Join(q, " AND ")
	return
}

func (s Service) searchSlsLog(ctx context.Context, chartsVisible bool, t1, t2 int64, pageSize, pageNo int, index, phrase, sql string) (items service.LoggerLogItems, count int, xAxis []string, yAxis []int64, err error) {
	phrase += " NOT 'sy-job-manager.39on.com'"
	res, num, err := s.LogStoreI.SearchLogs(ctx, t1, t2, index, phrase, sql, int64(pageSize), int64((pageNo-1)*pageSize))
	if err != nil {
		return
	}
	count = int(num)
	for _, v := range res {
		item, e := s.parseLogV2(v, index)
		if e != nil {
			continue
		}
		items = append(items, item)
	}
	if !chartsVisible {
		return
	}
	histograms, err := s.LogStoreI.GetSlsHistograms(ctx, t1, t2, index, phrase, sql)
	if err != nil {
		return
	}
	xAxis = make([]string, 0, len(histograms))
	yAxis = make([]int64, 0, len(histograms))
	for _, bucket := range histograms {
		xAxis = append(xAxis, gotil.DateSec(bucket.To))
		yAxis = append(yAxis, bucket.Count)
	}
	return
}

func (s Service) buildSearchQuery(ctx context.Context, req service.LoggerSearchRequest) (query elastic.Query) {
	return s.buildQuery(req.TimeA, req.TimeB, req.Or, req.Must, req.MustNot, req.WithoutOptions)
}

func (s Service) buildDateHistogramAggregation(ctx context.Context, req service.LoggerSearchRequest) *elastic.DateHistogramAggregation {
	interval := getInterval(req.TimeA, req.TimeB)
	return elastic.NewDateHistogramAggregation().
		Field(atTimestamp).
		FixedInterval(strconv.Itoa(interval) + "s").
		TimeZone("Asia/Shanghai").
		MinDocCount(0)
}

func getInterval(timeA, timeB int64) (interval int) {
	interval = (int)((timeB - timeA) / 1000 / lenChartPoints)
	if interval <= 1 {
		interval = 1
	} else if interval <= 5 {
		interval = 5
	} else if interval <= 10 {
		interval = 10
	} else if interval <= 30 {
		interval = 30
	} else if interval <= 60 {
		interval = 60
	} else if interval <= 300 {
		interval = 300
	} else if interval <= 900 {
		interval = 900
	} else if interval <= 1800 {
		interval = 1800
	} else if interval <= 3600 {
		interval = 3600
	} else if interval <= 3600*3 {
		interval = 3600 * 3
	} else if interval <= 3600*9 {
		interval = 3600 * 9
	} else if interval <= 3600*12 {
		interval = 3600 * 12
	} else {
		interval = 3600 * 24
	}
	return
}

type searchArgs struct {
	Client         *elastic.Client
	Index          string
	Query          elastic.Query
	Aggregation    elastic.Aggregation
	From           int
	Size           int
	Sort           string
	Raw            bool
	TrackTotalHits bool
}

func (s Service) search(ctx context.Context, args *searchArgs) (total int, items service.LoggerLogItems, xAxis []string, yAxis []int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	var (
		q      *elastic.SearchService
		result *elastic.SearchResult
	)
	q = args.Client.
		Search().
		Index(args.Index).
		TrackTotalHits(args.TrackTotalHits).
		Query(args.Query)
	if args.Aggregation != nil {
		q.Aggregation("charts", args.Aggregation)
	}
	result, err = q.
		From(args.From).
		Size(args.Size).
		Pretty(false).
		Version(true).
		Sort(args.Sort, false).
		Do(ctx)
	if err != nil {
		return
	}
	if args.TrackTotalHits {
		total = int(result.TotalHits())
	} else {
		total = 10000 // 返回一个数字即可，用于分页。
	}
	items = make(service.LoggerLogItems, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		item, e := s.parseLog(hit, args.Raw)
		if e != nil {
			continue
		}
		items = append(items, item)
	}
	sort.Sort(items)

	if args.Aggregation == nil {
		return
	}
	dateHistogram, ok := result.Aggregations.DateHistogram("charts")
	xAxis = make([]string, 0, len(dateHistogram.Buckets))
	yAxis = make([]int64, 0, len(dateHistogram.Buckets))
	if !ok {
		return
	}
	for _, bucket := range dateHistogram.Buckets {
		xAxis = append(xAxis, gotil.DateMs(int64(bucket.Key)))
		yAxis = append(yAxis, bucket.DocCount)
	}

	return
}
