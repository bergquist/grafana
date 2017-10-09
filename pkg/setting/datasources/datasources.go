package datasources

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"

	"github.com/grafana/grafana/pkg/models"
	yaml "gopkg.in/yaml.v2"
)

var (
	logger log.Logger
)

// TODO: secure jsonData
// TODO: auto reload on file changes

type DatasourcesAsConfig struct {
	PurgeOtherDatasources bool
	Datasources           []models.DataSource
}

func Init(configPath string) (error, io.Closer) {

	dc := newDatasourceConfiguration()
	dc.applyChanges(configPath)

	return nil, ioutil.NopCloser(nil)
}

type DatasourceConfigurator struct {
	log            log.Logger
	load           func(string) (*DatasourcesAsConfig, error)
	insert         func(*models.AddDataSourceCommand) error
	update         func(*models.UpdateDataSourceCommand) error
	delete         func(*models.DeleteDataSourceByIdCommand) error
	get            func(*models.GetDataSourceByNameQuery) error
	allDatasources func() ([]*models.DataSource, error)
}

func newDatasourceConfiguration() DatasourceConfigurator {
	return DatasourceConfigurator{
		log:    log.New("setting.datasource"),
		load:   readDatasources,
		insert: insertDatasource,
		update: updateDatasource,
		delete: deleteDatasource,
		get:    getDatasourceByName,
	}
}

func (dc *DatasourceConfigurator) applyChanges(configPath string) error {
	datasources, err := dc.load(configPath)
	if err != nil {
		return err
	}

	//read all datasources
	//delete datasources not in list

	for _, ds := range datasources.Datasources {
		if ds.OrgId == 0 {
			ds.OrgId = 1
		}

		query := &models.GetDataSourceByNameQuery{Name: ds.Name, OrgId: ds.OrgId}
		err := dc.get(query)
		if err != nil && err != models.ErrDataSourceNotFound {
			return err
		}

		fmt.Println("query", query.Result, query.Result.Name, query.Result.Id)
		if query.Result == nil {
			dc.log.Info("inserting ", "name", ds.Name)
			insertCmd := insertCommand(ds)
			if err := dc.insert(insertCmd); err != nil {
				return err
			}
		} else {
			dc.log.Info("updating", "name", ds.Name)
			updateCmd := updateCommand(ds, query.Result.Id)
			if err := dc.update(updateCmd); err != nil {
				return err
			}
		}
	}

	return nil
}

func readDatasources(path string) (*DatasourcesAsConfig, error) {
	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var datasources *DatasourcesAsConfig

	err = yaml.Unmarshal(yamlFile, &datasources)
	if err != nil {
		return nil, err
	}

	return datasources, nil
}

func insertCommand(ds models.DataSource) *models.AddDataSourceCommand {
	return &models.AddDataSourceCommand{
		OrgId:             ds.OrgId,
		Name:              ds.Name,
		Type:              ds.Type,
		Access:            ds.Access,
		Url:               ds.Url,
		Password:          ds.Password,
		User:              ds.User,
		Database:          ds.Database,
		BasicAuth:         ds.BasicAuth,
		BasicAuthUser:     ds.BasicAuthUser,
		BasicAuthPassword: ds.BasicAuthPassword,
		WithCredentials:   ds.WithCredentials,
		IsDefault:         ds.IsDefault,
		JsonData:          ds.JsonData,
	}
}

func updateCommand(ds models.DataSource, id int64) *models.UpdateDataSourceCommand {
	return &models.UpdateDataSourceCommand{
		Id:                id,
		OrgId:             ds.OrgId,
		Name:              ds.Name,
		Type:              ds.Type,
		Access:            ds.Access,
		Url:               ds.Url,
		Password:          ds.Password,
		User:              ds.User,
		Database:          ds.Database,
		BasicAuth:         ds.BasicAuth,
		BasicAuthUser:     ds.BasicAuthUser,
		BasicAuthPassword: ds.BasicAuthPassword,
		WithCredentials:   ds.WithCredentials,
		IsDefault:         ds.IsDefault,
		JsonData:          ds.JsonData,
	}
}

func deleteDatasource(cmd *models.DeleteDataSourceByIdCommand) error {
	return bus.Dispatch(cmd)
}

func updateDatasource(cmd *models.UpdateDataSourceCommand) error {
	return bus.Dispatch(cmd)
}

func insertDatasource(cmd *models.AddDataSourceCommand) error {
	return bus.Dispatch(cmd)
}

func getDatasourceByName(cmd *models.GetDataSourceByNameQuery) error {
	return bus.Dispatch(cmd)
}

func loadAllDatasources() ([]*models.DataSource, error) {
	dss := &models.GetAllDataSourcesQuery{}
	if err := bus.Dispatch(dss); err != nil {
		return nil, err
	}

	return dss.Result, nil
}
