package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
)

func TestLoggingList(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = []ari.LogData{
			ari.LogData{
				Key:    ari.NewKey(ari.LoggingKey, "c1"),
				Name:   "n1",
				Status: "s1",
				Types:  "t1",
			},
		}

		m.Logging.On("List", nil).Return(expected, nil)

		ld, err := cl.Asterisk().Logging().List(nil)
		if err != nil {
			t.Errorf("Unexpected error in logging list: %s", err)
		}
		if len(ld) != len(expected) {
			t.Errorf("Expected return of length %d, got %d", len(expected), len(ld))
		} else {
			for idx := range ld {
				failed := false
				failed = failed || ld[idx].ID != expected[idx].Name

				if failed {
					t.Errorf("Expected item '%d' to be '%v', got '%v",
						idx, expected[idx], ld[idx])
				}
			}
		}

		m.Logging.AssertCalled(t, "List", nil)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected []ari.LogData

		m.Logging.On("List").Return(expected, errors.New("error"))

		ld, err := cl.Asterisk().Logging().List(nil)
		if err == nil {
			t.Errorf("Expected error in logging list")
		}
		if len(ld) != len(expected) {
			t.Errorf("Expected return of length %d, got %d", len(expected), len(ld))
		} else {
			for idx := range ld {
				failed := false
				failed = failed || ld[idx].ID != expected[idx].Name

				if failed {
					t.Errorf("Expected item '%d' to be '%v', got '%v",
						idx, expected[idx], ld[idx])
				}
			}
		}

		m.Logging.AssertCalled(t, "List", nil)
	})
}

func TestLoggingCreate(t *testing.T, s Server) {
	key := ari.NewKey(ari.LoggingKey, "n1")
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Create", "n1", "l1").Return(nil)

		_, err := cl.Asterisk().Logging().Create(key, "l1")
		if err != nil {
			t.Errorf("Unexpected error in logging create: %s", err)
		}

		m.Logging.AssertCalled(t, "Create", key, "l1")
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Create", key, "l1").Return(errors.New("error"))

		_, err := cl.Asterisk().Logging().Create(key, "l1")
		if err == nil {
			t.Errorf("Expected error in logging create")
		}

		m.Logging.AssertCalled(t, "Create", key, "l1")
	})
}

func TestLoggingDelete(t *testing.T, s Server) {
	key := ari.NewKey(ari.LoggingKey, "n1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Delete", key).Return(nil)

		err := cl.Asterisk().Logging().Delete(key)
		if err != nil {
			t.Errorf("Unexpected error in logging Delete: %s", err)
		}

		m.Logging.AssertCalled(t, "Delete", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Delete", key).Return(errors.New("error"))

		err := cl.Asterisk().Logging().Delete(key)
		if err == nil {
			t.Errorf("Expected error in logging Delete")
		}

		m.Logging.AssertCalled(t, "Delete", key)
	})
}

func TestLoggingRotate(t *testing.T, s Server) {
	key := ari.NewKey(ari.LoggingKey, "n1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Rotate", key).Return(nil)

		err := cl.Asterisk().Logging().Rotate(key)
		if err != nil {
			t.Errorf("Unexpected error in logging Rotate: %s", err)
		}

		m.Logging.AssertCalled(t, "Rotate", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Logging.On("Rotate", key).Return(errors.New("error"))

		err := cl.Asterisk().Logging().Rotate(key)
		if err == nil {
			t.Errorf("Expected error in logging Rotate")
		}

		m.Logging.AssertCalled(t, "Rotate", key)
	})
}
