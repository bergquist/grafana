package tsdb

type QueryContext struct {
	TimeRange *TimeRange
	Queries   QuerySlice
	Results   map[string]*QueryResult
}

func NewQueryContext(queries QuerySlice, timeRange *TimeRange) *QueryContext {
	return &QueryContext{
		TimeRange: timeRange,
		Queries:   queries,
		Results:   make(map[string]*QueryResult),
	}
}
