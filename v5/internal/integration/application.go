package integration

import (
	"testing"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/pkg/errors"
	tmock "github.com/stretchr/testify/mock"
)

func TestApplicationList(t *testing.T, s Server) {
	runTest("emptyList", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Application.On("List", (*ari.Key)(nil)).Return([]*ari.Key{}, nil)

		if _, err := cl.Application().List(nil); err != nil {
			t.Errorf("Unexpected error in remote List call: %v", err)
		}
	})

	runTest("nonEmptyList", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		h1 := ari.NewKey(ari.ApplicationKey, "1")
		h2 := ari.NewKey(ari.ApplicationKey, "2")
		m.Application.On("List", (*ari.Key)(nil)).Return([]*ari.Key{h1, h2}, nil)

		items, err := cl.Application().List(nil)
		if err != nil {
			t.Errorf("Unexpected error in remote List call: %v", err)
		}
		if len(items) != 2 {
			t.Errorf("Expected items to be length 2, got %d", len(items))
		} else {
			if items[0].ID != "1" {
				t.Errorf("Expected item 0 to be '1', got %s", items[0].ID)
			}
			if items[1].ID != "2" {
				t.Errorf("Expected item 1 to be '2', got %s", items[1].ID)
			}
		}
	})
}

func TestApplicationData(t *testing.T, s Server) {
	key := ari.NewKey(ari.ApplicationKey, "1")
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		ad := &ari.ApplicationData{}
		ad.Name = "app1"

		m.Application.On("Data", tmock.Anything).Return(ad, nil)

		res, err := cl.Application().Data(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Data call: %v", err)
		}
		if res == nil || res.Name != ad.Name {
			t.Errorf("Expected application data name %s, got %s", ad, res)
		}

		m.Shutdown()
		m.Application.AssertCalled(t, "Data", key)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("unknown error")

		m.Application.On("Data", key).Return(nil, expected)

		res, err := cl.Application().Data(key)
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if res != nil {
			t.Errorf("Expected application data result to be empty, got %s", res)
		}

		m.Shutdown()

		m.Application.AssertCalled(t, "Data", key)
	})
}

func TestApplicationSubscribe(t *testing.T, s Server) {
	key := ari.NewKey(ari.ApplicationKey, "1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Application.On("Subscribe", key, "2").Return(nil)

		if err := cl.Application().Subscribe(key, "2"); err != nil {
			t.Errorf("Unexpected error in remote Subscribe call: %v", err)
		}

		m.Shutdown()

		m.Application.AssertCalled(t, "Subscribe", key, "2")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("unknown error")

		m.Application.On("Subscribe", key, "2").Return(expected)

		if err := cl.Application().Subscribe(key, "2"); err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}

		m.Shutdown()

		m.Application.AssertCalled(t, "Subscribe", key, "2")
	})
}

func TestApplicationUnsubscribe(t *testing.T, s Server) {
	key := ari.NewKey(ari.ApplicationKey, "1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Application.On("Unsubscribe", key, "2").Return(nil)

		if err := cl.Application().Unsubscribe(key, "2"); err != nil {
			t.Errorf("Unexpected error in remote Unsubscribe call: %T", err)
		}

		m.Shutdown()

		m.Application.AssertCalled(t, "Unsubscribe", key, "2")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("unknown error")

		m.Application.On("Unsubscribe", key, "2").Return(expected)

		if err := cl.Application().Unsubscribe(key, "2"); err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}

		m.Application.AssertCalled(t, "Unsubscribe", key, "2")
	})
}

func TestApplicationGet(t *testing.T, s Server) {
	key := ari.NewKey(ari.ApplicationKey, "1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		ad := &ari.ApplicationData{}
		ad.Name = "app1"

		m.Application.On("Data", tmock.Anything).Return(ad, nil)

		if h := cl.Application().Get(key); h == nil {
			t.Errorf("Unexpected nil-handle")
		}

		m.Shutdown()

		m.Application.AssertCalled(t, "Data", key)
	})
}
