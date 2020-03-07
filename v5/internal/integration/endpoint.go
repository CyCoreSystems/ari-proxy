package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari/v5"
)

func TestEndpointList(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		h1 := ari.NewEndpointKey("h1", "1")
		h2 := ari.NewEndpointKey("h1", "2")

		m.Endpoint.On("List", (*ari.Key)(nil)).Return([]*ari.Key{h1, h2}, nil)

		list, err := cl.Endpoint().List(nil)
		if err != nil {
			t.Errorf("Error in remote Endpoint List call: %s", err)
		}
		if len(list) != 2 {
			t.Errorf("Expected list of length 2, got %d", len(list))
		}

		m.Endpoint.AssertCalled(t, "List", (*ari.Key)(nil))
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Endpoint.On("List", (*ari.Key)(nil)).Return([]*ari.Key{}, errors.New("error"))

		list, err := cl.Endpoint().List(nil)
		if err == nil {
			t.Errorf("Expected error in remote Endpoint List call")
		}
		if len(list) != 0 {
			t.Errorf("Expected list of length 0, got %d", len(list))
		}

		m.Endpoint.AssertCalled(t, "List", (*ari.Key)(nil))
	})
}

func TestEndpointListByTech(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		h1 := ari.NewEndpointKey("h1", "1")
		h2 := ari.NewEndpointKey("h1", "2")

		m.Endpoint.On("ListByTech", "tech", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}).Return([]*ari.Key{h1, h2}, nil)

		list, err := cl.Endpoint().ListByTech("tech", nil)
		if err != nil {
			t.Errorf("Error in remote Endpoint List call: %s", err)
		}
		if len(list) != 2 {
			t.Errorf("Expected list of length 2, got %d", len(list))
		}

		m.Endpoint.AssertCalled(t, "ListByTech", "tech", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""})
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Endpoint.On("ListByTech", "tech", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}).Return([]*ari.Key{}, errors.New("error"))

		list, err := cl.Endpoint().ListByTech("tech", nil)
		if err == nil {
			t.Errorf("Expected error in remote Endpoint List call")
		}
		if len(list) != 0 {
			t.Errorf("Expected list of length 0, got %d", len(list))
		}

		m.Endpoint.AssertCalled(t, "ListByTech", "tech", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""})
	})
}

func TestEndpointData(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected ari.EndpointData
		expected.State = "st1"
		expected.Technology = "tech1"
		expected.Resource = "resource"

		h1 := ari.NewEndpointKey(expected.Technology, expected.Resource)

		m.Endpoint.On("Data", h1).Return(&expected, nil)

		data, err := cl.Endpoint().Data(h1)
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

		m.Endpoint.AssertCalled(t, "Data", h1)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected ari.EndpointData
		expected.State = "st1"
		expected.Technology = "tech1"
		expected.Resource = "resource"

		h1 := ari.NewEndpointKey(expected.Technology, expected.Resource)

		m.Endpoint.On("Data", h1).Return(nil, errors.New("error"))

		data, err := cl.Endpoint().Data(h1)
		if err == nil {
			t.Errorf("Expected error in remote Endpoint Data call")
		}
		if data != nil {
			t.Errorf("Expected data to be nil")
		}

		m.Endpoint.AssertCalled(t, "Data", h1)
	})
}
