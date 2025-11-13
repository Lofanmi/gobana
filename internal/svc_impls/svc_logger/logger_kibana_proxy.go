package svc_logger

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"unsafe"

	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/uritemplates"
)

func (s *Service) doSearchIndexKibanaProxy(ctx context.Context, cli *elastic.Client, req service.SearchRequest, indexName string, query elastic.Query, aggregation elastic.Aggregation) (result service.SearchResultElastic, err error) {
	result.RawQuery, result.ErrorResponse = query.Source()
	if result.ErrorResponse != nil {
		return
	}
	index := s.Indexes[indexName]
	indexList := index.Meta.IndexMetaForElastic.Indexes
	search := cli.Search()
	search.Index(indexList...).TrackTotalHits(req.TrackTotalHits).Query(query).Pretty(false).Version(true)
	search.From((req.PageNo - 1) * req.PageSize).Size(req.PageSize)
	if aggregation != nil {
		search.Aggregation("charts", aggregation)
	}
	for _, sortField := range index.Meta.IndexMetaForElastic.SortFields {
		search.Sort(sortField.Field, sortField.Ascending)
	}
	provider := s.Providers.Match(index.ProviderName)

	v := reflect.ValueOf(search).Elem()
	params := url.Values{}
	params.Set("method", http.MethodPost)
	field := v.FieldByName("index")
	if len(indexList) > 0 {
		path, _ := uritemplates.Expand("/{index}/_search", map[string]string{"index": strings.Join(indexList, ",")})
		params.Set("path", path)
	}
	field = v.FieldByName("searchSource")
	searchSource := (*elastic.SearchSource)(unsafe.Pointer(field.Pointer()))
	body, _ := searchSource.Source()
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Cookie", provider.AuthConfig.AuthForKibanaProxy.Cookie)
	headers.Set("kbn-version", provider.AuthConfig.AuthForKibanaProxy.KbnVersion)
	var res *elastic.Response
	if res, err = cli.PerformRequest(ctx, elastic.PerformRequestOptions{
		Method:  http.MethodPost,
		Path:    "/api/console/proxy",
		Params:  params,
		Body:    body,
		Headers: headers,
	}); err != nil {
		return
	}
	if err = json.Unmarshal(res.Body, &result.Response); err != nil {
		result.Response = new(elastic.SearchResult)
		result.Response.Header = res.Header
		result.ErrorResponse = err
		return
	}
	result.Response.Header = res.Header

	return
}
