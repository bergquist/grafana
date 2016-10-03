package elasticsearch

import (
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/tsdb"
)

type IndexPattern struct {
}

/*
IndexPattern.intervalMap = {
    "Hourly":   { startOf: 'hour',     amount: 'hours'},
    "Daily":   { startOf: 'day',      amount: 'days'},
    "Weekly":  { startOf: 'isoWeek',  amount: 'weeks'},
    "Monthly": { startOf: 'month',    amount: 'months'},
    "Yearly":  { startOf: 'year',     amount: 'years'},
  };
*/

var indexPattern map[string]IndexPattern = map[string]IndexPattern{
	"Hourly":  IndexPattern{},
	"Daily":   IndexPattern{},
	"Weekly":  IndexPattern{},
	"Monthly": IndexPattern{},
	"Yearly":  IndexPattern{},
}

type ElasticDatasource struct {
	*tsdb.DataSourceInfo

	ElasticVersion int
	TimeField      string
	Interval       string
	Index          string
}

func (datasource *ElasticDatasource) GetIndices(start, end time.Time) []string {
	index := datasource.Index

	var dates []time.Time
	var result []string
	dates = append(dates, start)

	for _, v := range dates {
		f := fmt.Sprintf("%v.%v.%v", v.Year(), int(v.Month()), v.Day())
		tmp := strings.Replace(strings.Replace(index, "[", "", 1), "]YYYY-MM-DD", f, 1)
		result = append(result, tmp)
	}

	return result
}
