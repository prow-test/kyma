package automock

import "k8s.io/helm/pkg/proto/hapi/chart"
import "github.com/kyma-project/kyma/components/helm-broker/internal"
import "github.com/stretchr/testify/mock"
import "github.com/Masterminds/semver"

// ChartStorage is an autogenerated mock type for the ChartStorage type
type ChartStorage struct {
	mock.Mock
}

// Get provides a mock function with given fields: name, ver
func (_m *ChartStorage) Get(name internal.ChartName, ver semver.Version) (*chart.Chart, error) {
	ret := _m.Called(name, ver)

	var r0 *chart.Chart
	if rf, ok := ret.Get(0).(func(internal.ChartName, semver.Version) *chart.Chart); ok {
		r0 = rf(name, ver)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chart.Chart)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(internal.ChartName, semver.Version) error); ok {
		r1 = rf(name, ver)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
