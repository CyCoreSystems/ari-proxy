package integration

import (
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
	"github.com/pkg/errors"
)

func TestApplicationList(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("emptyList", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Application.On("List").Return([]ari.ApplicationHandle{}, nil)

		if _, err := cl.Application().List(); err != nil {
			t.Errorf("Unexpected error in remote List call: %v", err)
		}
	})

	runTest("nonEmptyList", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		h1 := &mocks.ApplicationHandle{}
		h1.On("ID").Return("1")

		h2 := &mocks.ApplicationHandle{}
		h2.On("ID").Return("2")

		m.Application.On("List").Return([]ari.ApplicationHandle{h1, h2}, nil)

		items, err := cl.Application().List()
		if err != nil {
			t.Errorf("Unexpected error in remote List call: %v", err)
		}
		if len(items) != 2 {
			t.Errorf("Expected items to be length 2, got %d", len(items))
		} else {
			if items[0].ID() != "1" {
				t.Errorf("Expected item 0 to be '1', got %s", items[0].ID())
			}
			if items[1].ID() != "2" {
				t.Errorf("Expected item 1 to be '2', got %s", items[1].ID())
			}
		}

	})
}

func TestApplicationData(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		ad := &ari.ApplicationData{}
		ad.Name = "app1"

		m.Application.On("Data", "1").Return(ad, nil)

		res, err := cl.Application().Data("1")
		if err != nil {
			t.Errorf("Unexpected error in remote Data call: %v", err)
		}
		if res == nil || res.Name != ad.Name {
			t.Errorf("Expected application data name %s, got %s", ad, res)
		}
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("unknown error")

		m.Application.On("Data", "1").Return(nil, expected)

		res, err := cl.Application().Data("1")
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if res != nil {
			t.Errorf("Expected application data result to be empty, got %s", res)
		}
	})
}

func TestApplicationSubscribe(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Application.On("Subscribe", "1", "2").Return(nil)

		if err := cl.Application().Subscribe("1", "2"); err != nil {
			t.Errorf("Unexpected error in remote Subscribe call: %v", err)
		}

		m.Shutdown()

		m.Application.AssertCalled(t, "Subscribe", "1", "2")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("unknown error")

		m.Application.On("Subscribe", "1", "2").Return(expected)

		if err := cl.Application().Subscribe("1", "2"); err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}

		m.Shutdown()

		m.Application.AssertCalled(t, "Subscribe", "1", "2")
	})
}

func TestApplicationUnsubscribe(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Application.On("Unsubscribe", "1", "2").Return(nil)

		if err := cl.Application().Unsubscribe("1", "2"); err != nil {
			t.Errorf("Unexpected error in remote Unsubscribe call: %T", err)
		}

		m.Shutdown()

		m.Application.AssertCalled(t, "Unsubscribe", "1", "2")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("unknown error")

		m.Application.On("Unsubscribe", "1", "2").Return(expected)

		if err := cl.Application().Unsubscribe("1", "2"); err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}

		m.Application.AssertCalled(t, "Unsubscribe", "1", "2")
	})
}

func TestApplicationGet(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		if h := cl.Application().Get("1"); h == nil {
			t.Errorf("Unexpected nil-handle")
		}

		m.Shutdown()

		m.Application.AssertNotCalled(t, "Get", "1")
	})

}
