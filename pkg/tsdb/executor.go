package tsdb

import "context"
import "github.com/grafana/grafana/pkg/models"

type Executor interface {
	Execute(ctx context.Context, queries QuerySlice, query *TimeRange) *BatchResult
}

var registry map[string]GetExecutorFn

type GetExecutorFn func(dsInfo *models.DataSource) Executor

func init() {
	registry = make(map[string]GetExecutorFn)
}

func getExecutorFor(dsInfo *models.DataSource) Executor {
	if fn, exists := registry[dsInfo.Type]; exists {
		return fn(dsInfo)
	}
	return nil
}

func RegisterExecutor(pluginId string, fn GetExecutorFn) {
	registry[pluginId] = fn
}
