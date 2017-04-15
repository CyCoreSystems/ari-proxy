package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
)

func TestEndpointList(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		h1 := &mocks.EndpointHandle{}
		h2 := &mocks.EndpointHandle{}
		h1.On("ID").Return("h1/1")
		h2.On("ID").Return("h2/1")

		m.Endpoint.On("List").Return([]ari.EndpointHandle{h1, h2}, nil)

		list, err := cl.Endpoint().List()
		if err != nil {
			t.Errorf("Error in remote Endpoint List call: %s", err)
		}
		if len(list) != 2 {
			t.Errorf("Expected list of length 2, got %d", len(list))
		}

		h1.AssertCalled(t, "ID")
		h2.AssertCalled(t, "ID")
		m.Endpoint.AssertCalled(t, "List")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Endpoint.On("List").Return([]ari.EndpointHandle{}, errors.New("error"))

		list, err := cl.Endpoint().List()
		if err == nil {
			t.Errorf("Expected error in remote Endpoint List call")
		}
		if len(list) != 0 {
			t.Errorf("Expected list of length 0, got %d", len(list))
		}

		m.Endpoint.AssertCalled(t, "List")
	})
}

func TestEndpointListByTech(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		h1 := &mocks.EndpointHandle{}
		h2 := &mocks.EndpointHandle{}
		h1.On("ID").Return("h1/1")
		h2.On("ID").Return("h2/1")

		m.Endpoint.On("ListByTech", "tech").Return([]ari.EndpointHandle{h1, h2}, nil)

		list, err := cl.Endpoint().ListByTech("tech")
		if err != nil {
			t.Errorf("Error in remote Endpoint List call: %s", err)
		}
		if len(list) != 2 {
			t.Errorf("Expected list of length 2, got %d", len(list))
		}

		h1.AssertCalled(t, "ID")
		h2.AssertCalled(t, "ID")
		m.Endpoint.AssertCalled(t, "ListByTech", "tech")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Endpoint.On("ListByTech", "tech").Return([]ari.EndpointHandle{}, errors.New("error"))

		list, err := cl.Endpoint().ListByTech("tech")
		if err == nil {
			t.Errorf("Expected error in remote Endpoint List call")
		}
		if len(list) != 0 {
			t.Errorf("Expected list of length 0, got %d", len(list))
		}

		m.Endpoint.AssertCalled(t, "ListByTech", "tech")
	})
}

func TestEndpointData(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected ari.EndpointData
		expected.State = "st1"
		expected.Technology = "tech1"
		expected.Resource = "resource"

		m.Endpoint.On("Data", "tech", "resource").Return(&expected, nil)

		data, err := cl.Endpoint().Data("tech", "resource")
		if err != nil {
			t.Errorf("Error in remote Endpoint Data call: %s", err)
		}
		if data == nil {
			t.Errorf("Expected data to be non-nil")
		} else {
			failed := false
			failed = failed || expected.State != data.State
			failed = failed || expected.Resource != data.Resource
			failed = failed || expected.Technology != data.Technology
			if failed {
				t.Errorf("Expected '%v', got '%v'", expected, data)
			}
		}

		m.Endpoint.AssertCalled(t, "Data", "tech", "resource")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Endpoint.On("Data", "tech", "resource").Return(nil, errors.New("error"))

		data, err := cl.Endpoint().Data("tech", "resource")
		if err == nil {
			t.Errorf("Expected error in remote Endpoint Data call")
		}
		if data != nil {
			t.Errorf("Expected data to be nil")
		}

		m.Endpoint.AssertCalled(t, "Data", "tech", "resource")
	})
}
