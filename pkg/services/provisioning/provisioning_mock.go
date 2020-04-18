package provisioning

// Calls is a register of called functions
type Calls struct {
	ProvisionDatasources                []interface{}
	ProvisionNotifications              []interface{}
	ProvisionDashboards                 []interface{}
	GetDashboardProvisionerResolvedPath []interface{}
	GetAllowUIUpdatesFromConfig         []interface{}
}

// ServiceMock is a mock of `provisioning.Service`
type ServiceMock struct {
	Calls                                   *Calls
	ProvisionDatasourcesFunc                func() error
	ProvisionNotificationsFunc              func() error
	ProvisionDashboardsFunc                 func() error
	GetDashboardProvisionerResolvedPathFunc func(name string) string
	GetAllowUIUpdatesFromConfigFunc         func(name string) bool
}

// NewServiceMock returns a new provisioning service mock.
func NewServiceMock() *ServiceMock {
	return &ServiceMock{
		Calls: &Calls{},
	}
}

// ProvisionDatasources is a mock implementation
func (mock *ServiceMock) ProvisionDatasources() error {
	mock.Calls.ProvisionDatasources = append(mock.Calls.ProvisionDatasources, nil)
	if mock.ProvisionDatasourcesFunc != nil {
		return mock.ProvisionDatasourcesFunc()
	}
	return nil
}

// ProvisionNotifications is a mock implementation
func (mock *ServiceMock) ProvisionNotifications() error {
	mock.Calls.ProvisionNotifications = append(mock.Calls.ProvisionNotifications, nil)
	if mock.ProvisionNotificationsFunc != nil {
		return mock.ProvisionNotificationsFunc()
	}
	return nil
}

// ProvisionDashboards is a mock implementation
func (mock *ServiceMock) ProvisionDashboards() error {
	mock.Calls.ProvisionDashboards = append(mock.Calls.ProvisionDashboards, nil)
	if mock.ProvisionDashboardsFunc != nil {
		return mock.ProvisionDashboardsFunc()
	}
	return nil
}

// GetDashboardProvisionerResolvedPath is a mock implementation
func (mock *ServiceMock) GetDashboardProvisionerResolvedPath(name string) string {
	mock.Calls.GetDashboardProvisionerResolvedPath = append(mock.Calls.GetDashboardProvisionerResolvedPath, name)
	if mock.GetDashboardProvisionerResolvedPathFunc != nil {
		return mock.GetDashboardProvisionerResolvedPathFunc(name)
	}
	return ""
}

// GetAllowUIUpdatesFromConfig is a mock implementation
func (mock *ServiceMock) GetAllowUIUpdatesFromConfig(name string) bool {
	mock.Calls.GetAllowUIUpdatesFromConfig = append(mock.Calls.GetAllowUIUpdatesFromConfig, name)
	if mock.GetAllowUIUpdatesFromConfigFunc != nil {
		return mock.GetAllowUIUpdatesFromConfigFunc(name)
	}
	return false
}
