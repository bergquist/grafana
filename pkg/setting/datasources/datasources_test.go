package datasources

import (
	"fmt"
	"testing"

	"github.com/grafana/grafana/pkg/models"
)

func TestDatasourceAsConfig(t *testing.T) {
	fc := fakeConfiguration{
		cfg: &DatasourcesAsConfig{
			PurgeOtherDatasources: false,
			Datasources: []models.DataSource{
				models.DataSource{Name: "graphite", OrgId: 1},
			},
		},
		loadAll: []*models.DataSource{
			&models.DataSource{Name: "graphite", OrgId: 1, Id: 1},
		},
	}

	dc := fakeConfigurator(fc)

	err := dc.applyChanges("")
	if err != nil {
		t.Fatalf("applyChanges return an error %v", err)
	}

	if len(fc.updated) != 1 {

	}
}

func fakeConfigurator(fc fakeConfiguration) DatasourceConfigurator {
	dc := newDatasourceConfiguration()
	dc.insert = fc.insertDatasource
	dc.update = fc.updateDatasource
	dc.delete = fc.deleteDatasource
	dc.allDatasources = fc.loadAllDatasources
	dc.load = fc.load
	dc.get = fc.getDatasourceByName

	return dc
}

type fakeConfiguration struct {
	inserted []*models.AddDataSourceCommand
	deleted  []*models.DeleteDataSourceByIdCommand
	updated  []*models.UpdateDataSourceCommand

	loadAll []*models.DataSource

	cfg *DatasourcesAsConfig
}

func (fc *fakeConfiguration) load(path string) (*DatasourcesAsConfig, error) {
	return fc.cfg, nil
}

func (fc *fakeConfiguration) deleteDatasource(cmd *models.DeleteDataSourceByIdCommand) error {
	fc.deleted = append(fc.deleted, cmd)
	return nil
}

func (fc *fakeConfiguration) updateDatasource(cmd *models.UpdateDataSourceCommand) error {
	fc.updated = append(fc.updated, cmd)
	return nil
}

func (fc *fakeConfiguration) insertDatasource(cmd *models.AddDataSourceCommand) error {
	fc.inserted = append(fc.inserted, cmd)
	return nil
}

func (fc *fakeConfiguration) loadAllDatasources() ([]*models.DataSource, error) {
	return fc.loadAll, nil
}

func (fc *fakeConfiguration) getDatasourceByName(cmd *models.GetDataSourceByNameQuery) error {
	for _, v := range fc.loadAll {
		fmt.Printf("cmd %v %v \n v %v %v", cmd.Name, cmd.OrgId, v.Name, v.OrgId)
		if cmd.Name == v.Name && cmd.OrgId == v.OrgId {

			ds := *v
			cmd.Result = &ds
			fmt.Printf("returnning %v\n", cmd.Result)
			return nil
		}
	}

	return models.ErrDataSourceNotFound
}
