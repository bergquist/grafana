package tsdb

import (
	"context"

	"fmt"

	"github.com/grafana/grafana/pkg/models"
)

type Executor interface {
	Execute(ctx context.Context, queries QuerySlice, query *TimeRange) *BatchResult
}

var registry map[string]GetExecutorFn

type GetExecutorFn func(dsInfo *models.DataSource) (Executor, error)

func init() {
	registry = make(map[string]GetExecutorFn)
}

func getExecutorFor(dsInfo *models.DataSource) (Executor, error) {
	if fn, exists := registry[dsInfo.Type]; exists {
		return fn(dsInfo)
	}
	return nil, fmt.Errorf("Could not find executor for datasource: %s", dsInfo.Type)
}

func RegisterExecutor(pluginId string, fn GetExecutorFn) {
	registry[pluginId] = fn
}
