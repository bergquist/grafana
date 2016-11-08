package influxdb

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/grafana/pkg/tsdb"
	"gopkg.in/guregu/null.v3"
)

type ResponseParser struct{}

func (rp *ResponseParser) Parse(response *Response, query *Query) *tsdb.QueryResult {
	queryRes := tsdb.NewQueryResult()

	for _, result := range response.Results {
		queryRes.Series = append(queryRes.Series, rp.transformRows(result.Series, queryRes, query)...)
	}

	return queryRes
}

func (rp *ResponseParser) transformRows(rows []Row, queryResult *tsdb.QueryResult, query *Query) tsdb.TimeSeriesSlice {
	var result tsdb.TimeSeriesSlice

	for _, row := range rows {
		for columnIndex, column := range row.Columns {
			if column == "time" {
				continue
			}

			var points tsdb.TimeSeriesPoints
			for _, valuePair := range row.Values {
				point, err := rp.parseTimepoint(valuePair, columnIndex)
				if err == nil {
					points = append(points, point)
				}
			}
			result = append(result, &tsdb.TimeSeries{
				Name:   rp.formatSerieName(row, column, query),
				Points: points,
			})
		}
	}

	return result
}

func (rp *ResponseParser) formatSerieName(row Row, column string, query *Query) string {
	if query.Alias == "" {
		return rp.buildSerieNameFromQuery(row, column)
	}

	reg, _ := regexp.Compile(`\$\s*(.+?)\s|`)

	result := reg.ReplaceAllFunc([]byte(query.Alias), func(in []byte) []byte {
		sin := string(in)

		if strings.HasPrefix(sin, "m") || strings.HasPrefix(sin, "measurement") {
			return []byte(query.Measurement)
		}
		if strings.HasPrefix(sin, "col") {
			return []byte(column)
		}

		if !strings.HasPrefix(sin, "tag_") {
			return in
		}

		//\$\s*(.+?)\s|
		//|\[\[([\s*(.+?)*\s*]+?)\]\]

		//labelName := strings.Replace(string(in), "{{", "", 1)
		//labelName = strings.Replace(labelName, "}}", "", 1)
		//labelName = strings.TrimSpace(labelName)
		//if val, exists := metric[pmodel.LabelName(labelName)]; exists {
		//	return []byte(val)
		//}

		return in
	})

	return string(result)
}

/*
   var regex = /\$(\w+)|\[\[([\s\S]+?)\]\]/g;
   var segments = series.name.split('.');

   return this.alias.replace(regex, function(match, g1, g2) {
     var group = g1 || g2;
     var segIndex = parseInt(group, 10);

     if (group === 'm' || group === 'measurement') { return series.name; }
     if (group === 'col') { return series.columns[index]; }
     if (!isNaN(segIndex)) { return segments[segIndex]; }
     if (group.indexOf('tag_') !== 0) { return match; }

     var tag = group.replace('tag_', '');
     if (!series.tags) { return match; }
     return series.tags[tag];
   });
*/

func (rp *ResponseParser) buildSerieNameFromQuery(row Row, column string) string {
	var tags []string

	for k, v := range row.Tags {
		tags = append(tags, fmt.Sprintf("%s: %s", k, v))
	}

	tagText := ""
	if len(tags) > 0 {
		tagText = fmt.Sprintf(" { %s }", strings.Join(tags, " "))
	}

	return fmt.Sprintf("%s.%s%s", row.Name, column, tagText)
}

func (rp *ResponseParser) parseTimepoint(valuePair []interface{}, valuePosition int) (tsdb.TimePoint, error) {
	var value null.Float = rp.parseValue(valuePair[valuePosition])

	timestampNumber, _ := valuePair[0].(json.Number)
	timestamp, err := timestampNumber.Float64()
	if err != nil {
		return tsdb.TimePoint{}, err
	}

	return tsdb.NewTimePoint(value, timestamp), nil
}

func (rp *ResponseParser) parseValue(value interface{}) null.Float {
	number, ok := value.(json.Number)
	if !ok {
		return null.FloatFromPtr(nil)
	}

	fvalue, err := number.Float64()
	if err == nil {
		return null.FloatFrom(fvalue)
	}

	ivalue, err := number.Int64()
	if err == nil {
		return null.FloatFrom(float64(ivalue))
	}

	return null.FloatFromPtr(nil)
}
