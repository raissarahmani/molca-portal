package analytics

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/molca-id/portal-app-api/api/project/model"
	coreAnalytics "github.com/molca-id/portal-app-api/arch/analytics"
	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	// analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

type RankingItem struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type TimeSeriesItem struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

type Service interface {
	GetVisitsTimeSeries(rangeType string, debug bool) ([]TimeSeriesItem, error)
	GetRankingByTitle(order string, limit int64, rangeType string, debug bool) ([]RankingItem, error)
	GetProjectVisitByType(order string, limit int64, rangeType string, debug bool) ([]RankingItem, error)
	GetProjectType(debug bool) ([]RankingItem, error)
}

type service struct {
	network.BaseService
	coreAnalytics       *coreAnalytics.Service
	projectQueryBuilder mongo.QueryBuilder[model.Project]
	db                  mongo.Database
}

func NewService(coreAnalytics *coreAnalytics.Service, db mongo.Database) Service {
	return &service{
		BaseService:         network.NewBaseService(),
		coreAnalytics:       coreAnalytics,
		projectQueryBuilder: mongo.NewQueryBuilder[model.Project](db, model.CollectionName),
		db:                  db,
	}
}

/*
GetVisitsTimeSeries:
- Always fetch daily rows from arch (safer to aggregate locally).
- Parse GA rows (YYYYMMDD) -> TimeSeriesItem{Date:"2006-01-02", Value: int64}
- If rangeType == weekly -> group into "MonShort Week N" (e.g. "Sep Week 2") summing values.
- If rangeType == monthly -> group into "Jan" style month labels summing values.
- Default (daily) -> format date to "02/01" (dd/mm) and return.
*/
func (s *service) GetVisitsTimeSeries(rangeType string, debug bool) ([]TimeSeriesItem, error) {
	resp, err := s.coreAnalytics.GetVisitsTimeSeries("daily", debug)
	if err != nil {
		return nil, fmt.Errorf("failed to get user visits by range: %w", err)
	}

	var daily []TimeSeriesItem
	for _, row := range resp.Rows {
		if len(row.DimensionValues) == 0 || len(row.MetricValues) == 0 {
			continue
		}
		v, err := strconv.ParseInt(row.MetricValues[0].Value, 10, 64)
		if err != nil {
			continue
		}
		rawDate := row.DimensionValues[0].Value
		parsed, err := time.Parse("20060102", rawDate)
		if err != nil {
			continue
		}
		daily = append(daily, TimeSeriesItem{
			Date:  parsed.Format("2006-01-02"),
			Value: v,
		})
	}

	switch rangeType {
	case "weekly":
		return groupByWeek(daily), nil
	case "monthly":
		return groupByMonth(daily), nil
	default:
		return formatDaily(daily), nil
	}
}

/*
GetRankingByTitle:
- For now returns mocked top-10 project titles with counts.
- Supports order (asc|desc) and limit (applied after sorting).
- Accepts rangeType & debug but mocks ignore them; later replace mock with real GA call.
*/
func (s *service) GetRankingByTitle(order string, limit int64, rangeType string, debug bool) ([]RankingItem, error) {
	mock := []RankingItem{
		{Value: "dg-1", Count: 12},
		{Value: "ar-1", Count: 18},
		{Value: "vr-1", Count: 16},
		{Value: "sm-1", Count: 14},
		{Value: "deck-1", Count: 15},
		{Value: "tool-1", Count: 10},
		{Value: "dg-2", Count: 12},
		{Value: "sm-2", Count: 10},
		{Value: "deck-2", Count: 10},
		{Value: "tool-2", Count: 20},
	}

	if limit <= 0 {
		limit = 6
	}

	if order == "asc" {
		sort.Slice(mock, func(i, j int) bool { return mock[i].Count < mock[j].Count })
	} else {
		sort.Slice(mock, func(i, j int) bool { return mock[i].Count > mock[j].Count })
	}

	if int64(len(mock)) > limit {
		mock = mock[:limit]
	}
	return mock, nil
}

/*
GetProjectVisitByType:
- Mocked 6 types (exactly). No sorting requested: we keep a deterministic order.
- Accepts order/limit/rangeType/debug for compatibility but mock ignores range aggregation.
*/
func (s *service) GetProjectVisitByType(order string, limit int64, rangeType string, debug bool) ([]RankingItem, error) {
	mock := []RankingItem{
		{Value: "digital-twin", Count: 20},
		{Value: "vr", Count: 10},
		{Value: "ar", Count: 15},
		{Value: "smart-manufacture", Count: 20},
		{Value: "deck", Count: 15},
		{Value: "tool", Count: 15},
	}

	if limit > 0 && int64(len(mock)) > limit {
		mock = mock[:limit]
	}
	return mock, nil
}

/*
GetProjectType:
- Pull projects from DB (projection only type), count occurrences per type (6 types).
- Return counts sorted descending (so frontend gets biggest first).
*/
func (s *service) GetProjectType(debug bool) ([]RankingItem, error) {
	opts := options.Find().SetProjection(bson.M{"type": 1})
	projects, err := s.projectQueryBuilder.
		SingleQuery().
		FindAll(bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	typeCount := make(map[string]int64)
	for _, p := range projects {
		if p.Type != "" {
			typeCount[p.Type]++
		}
	}

	var ranking []RankingItem
	for t, c := range typeCount {
		ranking = append(ranking, RankingItem{Value: t, Count: c})
	}

	sort.Slice(ranking, func(i, j int) bool { return ranking[i].Count > ranking[j].Count })
	return ranking, nil
}

/* parseRankingResponse retained for future GA integration */
// func (s *service) parseRankingResponse(resp *analyticsdata.RunReportResponse) []RankingItem {
// 	var ranking []RankingItem
// 	for _, row := range resp.Rows {
// 		if len(row.DimensionValues) > 0 && len(row.MetricValues) > 0 {
// 			count, err := strconv.ParseInt(row.MetricValues[0].Value, 10, 64)
// 			if err != nil {
// 				continue
// 			}
// 			ranking = append(ranking, RankingItem{
// 				Value: row.DimensionValues[0].Value,
// 				Count: count,
// 			})
// 		}
// 	}
// 	return ranking
// }

func formatDaily(items []TimeSeriesItem) []TimeSeriesItem {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Date < items[j].Date
	})
	for i := range items {
		t, err := time.Parse("2006-01-02", items[i].Date)
		if err == nil {
			items[i].Date = t.Format("02/01")
		}
	}
	return items
}

func groupByWeek(items []TimeSeriesItem) []TimeSeriesItem {
	buckets := make(map[string]int64)
	labels := make(map[string]string)

	for _, it := range items {
		t, err := time.Parse("2006-01-02", it.Date)
		if err != nil {
			continue
		}
		year := t.Year()
		month := int(t.Month())
		weekOfMonth := (t.Day()-1)/7 + 1
		key := fmt.Sprintf("%04d-%02d-%d", year, month, weekOfMonth)
		monthShort := t.Format("Jan")
		buckets[key] += it.Value
		labels[key] = fmt.Sprintf("%s Week %d", monthShort, weekOfMonth)
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	res := make([]TimeSeriesItem, 0, len(keys))
	for _, k := range keys {
		res = append(res, TimeSeriesItem{Date: labels[k], Value: buckets[k]})
	}
	return res
}

func groupByMonth(items []TimeSeriesItem) []TimeSeriesItem {
	buckets := make(map[string]int64)
	labels := make(map[string]string)

	for _, it := range items {
		t, err := time.Parse("2006-01-02", it.Date)
		if err != nil {
			continue
		}
		year := t.Year()
		month := int(t.Month())
		key := fmt.Sprintf("%04d-%02d", year, month)
		buckets[key] += it.Value
		labels[key] = t.Format("Jan")
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	res := make([]TimeSeriesItem, 0, len(keys))
	for _, k := range keys {
		res = append(res, TimeSeriesItem{Date: labels[k], Value: buckets[k]})
	}
	return res
}
