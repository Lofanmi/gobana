package svc_logger

import (
	"github.com/Lofanmi/cmap"
	"github.com/Lofanmi/gobana/service"
	"github.com/spf13/cast"
)

func (s *Service) ParseTotal(m *cmap.Map[string, service.SearchResult]) (total int) {
	indexes := m.Keys()
	for _, index := range indexes {
		searchResult, exist := m.Get(index)
		if !exist {
			continue
		}
		switch searchResult.ProviderType {
		case service.ProviderTypeSls:
			s.parseTotalSls(searchResult.SearchResultSls, &total)
		case service.ProviderTypeElasticsearch:
			s.parseTotalElastic(searchResult.SearchResultElastic, &total)
		}
	}
	return
}

func (s *Service) parseTotalSls(sr service.SearchResultSls, total *int) {
	if sr.ResponseLog == nil || len(sr.ResponseLog.Logs) <= 0 {
		return
	}
	if sr.ResponseCountByGetHistograms != nil {
		*total = *total + int(sr.ResponseCountByGetHistograms.Count)
		return
	}
	if sr.ResponseCountByGetLogs != nil && len(sr.ResponseCountByGetLogs.Logs) > 0 {
		log := sr.ResponseCountByGetLogs.Logs[0]
		if log == nil {
			return
		}
		*total = *total + cast.ToInt(log["count"])
	}
}

func (s *Service) parseTotalElastic(sr service.SearchResultElastic, total *int) {
	result := sr.Response
	if result == nil || result.Hits == nil || result.Hits.TotalHits == nil {
		return
	}
	*total = *total + int(result.Hits.TotalHits.Value)
}
