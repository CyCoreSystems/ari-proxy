package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
)

func TestLoggingList(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = []ari.LogData{
			ari.LogData{
				Configuration: "c1",
				Name:          "n1",
				Status:        "s1",
				Type:          "t1",
			},
		}

		m.Logging.On("List").Return(expected, nil)

		ld, err := cl.Asterisk().Logging().List()
		if err != nil {
			t.Errorf("Unexpected error in logging list: %s", err)
		}
		if len(ld) != len(expected) {
			t.Errorf("Expected return of length %d, got %d", len(expected), len(ld))
		} else {
			for idx := range ld {
				failed := false
				failed = failed || ld[idx].Configuration != expected[idx].Configuration
				failed = failed || ld[idx].Name != expected[idx].Name
				failed = failed || ld[idx].Status != expected[idx].Status
				failed = failed || ld[idx].Type != expected[idx].Type

				if failed {
					t.Errorf("Expected item '%d' to be '%v', got '%v",
						idx, expected[idx], ld[idx])
				}
			}
		}

		m.Logging.AssertCalled(t, "List")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var expected []ari.LogData

		m.Logging.On("List").Return(expected, errors.New("error"))

		ld, err := cl.Asterisk().Logging().List()
		if err == nil {
			t.Errorf("Expected error in logging list")
		}
		if len(ld) != len(expected) {
			t.Errorf("Expected return of length %d, got %d", len(expected), len(ld))
		} else {
			for idx := range ld {
				failed := false
				failed = failed || ld[idx].Configuration != expected[idx].Configuration
				failed = failed || ld[idx].Name != expected[idx].Name
				failed = failed || ld[idx].Status != expected[idx].Status
				failed = failed || ld[idx].Type != expected[idx].Type

				if failed {
					t.Errorf("Expected item '%d' to be '%v', got '%v",
						idx, expected[idx], ld[idx])
				}
			}
		}

		m.Logging.AssertCalled(t, "List")
	})
}

func TestLoggingCreate(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Create", "n1", "l1").Return(nil)

		err := cl.Asterisk().Logging().Create("n1", "l1")
		if err != nil {
			t.Errorf("Unexpected error in logging create: %s", err)
		}

		m.Logging.AssertCalled(t, "Create", "n1", "l1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Create", "n1", "l1").Return(errors.New("error"))

		err := cl.Asterisk().Logging().Create("n1", "l1")
		if err == nil {
			t.Errorf("Expected error in logging create")
		}

		m.Logging.AssertCalled(t, "Create", "n1", "l1")
	})
}

func TestLoggingDelete(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Delete", "n1").Return(nil)

		err := cl.Asterisk().Logging().Delete("n1")
		if err != nil {
			t.Errorf("Unexpected error in logging Delete: %s", err)
		}

		m.Logging.AssertCalled(t, "Delete", "n1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Delete", "n1").Return(errors.New("error"))

		err := cl.Asterisk().Logging().Delete("n1")
		if err == nil {
			t.Errorf("Expected error in logging Delete")
		}

		m.Logging.AssertCalled(t, "Delete", "n1")
	})
}

func TestLoggingRotate(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Rotate", "n1").Return(nil)

		err := cl.Asterisk().Logging().Rotate("n1")
		if err != nil {
			t.Errorf("Unexpected error in logging Rotate: %s", err)
		}

		m.Logging.AssertCalled(t, "Rotate", "n1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Rotate", "n1").Return(errors.New("error"))

		err := cl.Asterisk().Logging().Rotate("n1")
		if err == nil {
			t.Errorf("Expected error in logging Rotate")
		}

		m.Logging.AssertCalled(t, "Rotate", "n1")
	})
}
