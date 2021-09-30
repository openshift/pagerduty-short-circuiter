// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/client/client.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	http "net/http"
	reflect "reflect"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	gomock "github.com/golang/mock/gomock"
)

// MockPagerDutyClient is a mock of PagerDutyClient interface.
type MockPagerDutyClient struct {
	ctrl     *gomock.Controller
	recorder *MockPagerDutyClientMockRecorder
}

// MockPagerDutyClientMockRecorder is the mock recorder for MockPagerDutyClient.
type MockPagerDutyClientMockRecorder struct {
	mock *MockPagerDutyClient
}

// NewMockPagerDutyClient creates a new mock instance.
func NewMockPagerDutyClient(ctrl *gomock.Controller) *MockPagerDutyClient {
	mock := &MockPagerDutyClient{ctrl: ctrl}
	mock.recorder = &MockPagerDutyClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPagerDutyClient) EXPECT() *MockPagerDutyClientMockRecorder {
	return m.recorder
}

// GetCurrentUser mocks base method.
func (m *MockPagerDutyClient) GetCurrentUser(arg0 pagerduty.GetCurrentUserOptions) (*pagerduty.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentUser", arg0)
	ret0, _ := ret[0].(*pagerduty.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentUser indicates an expected call of GetCurrentUser.
func (mr *MockPagerDutyClientMockRecorder) GetCurrentUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentUser", reflect.TypeOf((*MockPagerDutyClient)(nil).GetCurrentUser), arg0)
}

// GetIncidentAlert mocks base method.
func (m *MockPagerDutyClient) GetIncidentAlert(incidentID, alertID string) (*pagerduty.IncidentAlertResponse, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIncidentAlert", incidentID, alertID)
	ret0, _ := ret[0].(*pagerduty.IncidentAlertResponse)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetIncidentAlert indicates an expected call of GetIncidentAlert.
func (mr *MockPagerDutyClientMockRecorder) GetIncidentAlert(incidentID, alertID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIncidentAlert", reflect.TypeOf((*MockPagerDutyClient)(nil).GetIncidentAlert), incidentID, alertID)
}

// GetService mocks base method.
func (m *MockPagerDutyClient) GetService(serviceID string, opts *pagerduty.GetServiceOptions) (*pagerduty.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetService", serviceID, opts)
	ret0, _ := ret[0].(*pagerduty.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetService indicates an expected call of GetService.
func (mr *MockPagerDutyClientMockRecorder) GetService(serviceID, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetService", reflect.TypeOf((*MockPagerDutyClient)(nil).GetService), serviceID, opts)
}

// ListIncidentAlerts mocks base method.
func (m *MockPagerDutyClient) ListIncidentAlerts(incidentId string) (*pagerduty.ListAlertsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIncidentAlerts", incidentId)
	ret0, _ := ret[0].(*pagerduty.ListAlertsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIncidentAlerts indicates an expected call of ListIncidentAlerts.
func (mr *MockPagerDutyClientMockRecorder) ListIncidentAlerts(incidentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIncidentAlerts", reflect.TypeOf((*MockPagerDutyClient)(nil).ListIncidentAlerts), incidentId)
}

// ListIncidents mocks base method.
func (m *MockPagerDutyClient) ListIncidents(arg0 pagerduty.ListIncidentsOptions) (*pagerduty.ListIncidentsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIncidents", arg0)
	ret0, _ := ret[0].(*pagerduty.ListIncidentsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIncidents indicates an expected call of ListIncidents.
func (mr *MockPagerDutyClientMockRecorder) ListIncidents(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIncidents", reflect.TypeOf((*MockPagerDutyClient)(nil).ListIncidents), arg0)
}
