package analytics

import (
	"context"
	"fmt"
	"strings"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

type Service struct {
	analytics  *analyticsdata.Service
	propertyID string
}

func NewService(ctx context.Context, propertyID string, credentialsPath string) (*Service, error) {
	analyticsService, err := analyticsdata.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, fmt.Errorf("unable to create analytics data service: %w", err)
	}

	return &Service{
		analytics:  analyticsService,
		propertyID: propertyID,
	}, nil
}

func mapRangeToDimensionAndDate(rangeType string) (string, string) {
	switch rangeType {
	case "daily":
		return "28daysAgo", "date"
	case "weekly":
		return "12weeksAgo", "week"
	case "monthly":
		return "12monthsAgo", "month"
	default:
		return "28daysAgo", "date"
	}
}

func buildDebugFilter() *analyticsdata.FilterExpression {
	return &analyticsdata.FilterExpression{
		Filter: &analyticsdata.Filter{
			FieldName: "debug_mode",
			StringFilter: &analyticsdata.StringFilter{
				MatchType: "EXACT",
				Value:     "true",
			},
		},
	}
}

func (s *Service) GetVisitsTimeSeries(rangeType string, debug bool) (*analyticsdata.RunReportResponse, error) {
	startDate, dateDimension := mapRangeToDimensionAndDate(rangeType)

	req := &analyticsdata.RunReportRequest{
		Dimensions: []*analyticsdata.Dimension{
			{Name: dateDimension},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "screenPageViews"},
		},
		DateRanges: []*analyticsdata.DateRange{
			{
				StartDate: startDate,
				EndDate:   "today",
			},
		},
		OrderBys: []*analyticsdata.OrderBy{
			{
				Dimension: &analyticsdata.DimensionOrderBy{
					DimensionName: dateDimension,
				},
				Desc: false,
			},
		},
	}

	if debug {
		req.DimensionFilter = buildDebugFilter()
	}

	return s.analytics.Properties.RunReport("properties/"+s.propertyID, req).Do()
}

func (s *Service) GetProjectClicksByTitle(rangeType string, order string, debug bool) (*analyticsdata.RunReportResponse, error) {
	startDate, dateDimension := mapRangeToDimensionAndDate(rangeType)
	desc := true
	if strings.ToLower(order) == "asc" {
		desc = false
	}

	req := &analyticsdata.RunReportRequest{
		Dimensions: []*analyticsdata.Dimension{
			{Name: dateDimension},
			{Name: "customEvent:project_title"},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "eventCount"},
		},
		DateRanges: []*analyticsdata.DateRange{
			{
				StartDate: startDate,
				EndDate:   "today",
			},
		},
		OrderBys: []*analyticsdata.OrderBy{
			{
				Dimension: &analyticsdata.DimensionOrderBy{
					DimensionName: dateDimension,
				},
				Desc: true,
			},
			{
				Metric: &analyticsdata.MetricOrderBy{
					MetricName: "eventCount",
				},
				Desc: desc,
			},
		},
		DimensionFilter: &analyticsdata.FilterExpression{
			Filter: &analyticsdata.Filter{
				FieldName: "eventName",
				StringFilter: &analyticsdata.StringFilter{
					MatchType: "EXACT",
					Value:     "project_click",
				},
			},
		},
	}

	if debug {
		req.DimensionFilter = &analyticsdata.FilterExpression{
			AndGroup: &analyticsdata.FilterExpressionList{
				Expressions: []*analyticsdata.FilterExpression{
					req.DimensionFilter,
					buildDebugFilter(),
				},
			},
		}
	}

	return s.analytics.Properties.RunReport("properties/"+s.propertyID, req).Do()
}

func (s *Service) GetProjectClicksByType(rangeType string, debug bool) (*analyticsdata.RunReportResponse, error) {
	startDate, dateDimension := mapRangeToDimensionAndDate(rangeType)

	req := &analyticsdata.RunReportRequest{
		Dimensions: []*analyticsdata.Dimension{
			{Name: dateDimension},
			{Name: "customEvent:project_type"},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "eventCount"},
		},
		DateRanges: []*analyticsdata.DateRange{
			{
				StartDate: startDate,
				EndDate:   "today",
			},
		},
		OrderBys: []*analyticsdata.OrderBy{
			{
				Dimension: &analyticsdata.DimensionOrderBy{
					DimensionName: dateDimension,
				},
				Desc: true,
			},
		},
		DimensionFilter: &analyticsdata.FilterExpression{
			Filter: &analyticsdata.Filter{
				FieldName: "eventName",
				StringFilter: &analyticsdata.StringFilter{
					MatchType: "EXACT",
					Value:     "project_click",
				},
			},
		},
	}

	if debug {
		req.DimensionFilter = &analyticsdata.FilterExpression{
			AndGroup: &analyticsdata.FilterExpressionList{
				Expressions: []*analyticsdata.FilterExpression{
					req.DimensionFilter,
					buildDebugFilter(),
				},
			},
		}
	}

	return s.analytics.Properties.RunReport("properties/"+s.propertyID, req).Do()
}
