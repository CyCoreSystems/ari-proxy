package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari/v5"
)

func TestMailboxList(t *testing.T, s Server) {
	runTest("empty", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Mailbox.On("List", (*ari.Key)(nil)).Return([]*ari.Key{}, nil)

		ret, err := cl.Mailbox().List(nil)
		if err != nil {
			t.Errorf("Unexpected error in remote List call")
		}
		if len(ret) != 0 {
			t.Errorf("Expected return length to be 0, got %d", len(ret))
		}

		m.Shutdown()

		m.Mailbox.AssertCalled(t, "List", (*ari.Key)(nil))
	})

	runTest("nonEmpty", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		h1 := ari.NewKey(ari.MailboxKey, "h1")
		h2 := ari.NewKey(ari.MailboxKey, "h2")

		m.Mailbox.On("List", (*ari.Key)(nil)).Return([]*ari.Key{h1, h2}, nil)

		ret, err := cl.Mailbox().List(nil)
		if err != nil {
			t.Errorf("Unexpected error in remote List call")
		}
		if len(ret) != 2 {
			t.Errorf("Expected return length to be 2, got %d", len(ret))
		}

		m.Shutdown()

		m.Mailbox.AssertCalled(t, "List", (*ari.Key)(nil))
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Mailbox.On("List", (*ari.Key)(nil)).Return(nil, errors.New("unknown error"))

		ret, err := cl.Mailbox().List(nil)
		if err == nil {
			t.Errorf("Expected error in remote List call")
		}
		if len(ret) != 0 {
			t.Errorf("Expected return length to be 0, got %d", len(ret))
		}

		m.Shutdown()

		m.Mailbox.AssertCalled(t, "List", (*ari.Key)(nil))
	})
}

func testMailboxCommand(t *testing.T, m *mock, name string, id *ari.Key, expected error, fn func(*ari.Key) error) {
	m.Mailbox.On(name, id).Return(expected)
	err := fn(id)
	failed := false
	failed = failed || err == nil && expected != nil
	failed = failed || err != nil && expected == nil
	failed = failed || err != nil && expected != nil && err.Error() != expected.Error()
	if failed {
		t.Errorf("Expected mailbox %s(%s) to return '%v', got '%v'",
			name, id, expected, err,
		)
	}
	m.Mailbox.AssertCalled(t, name, id)
}

func TestMailboxDelete(t *testing.T, s Server) {
	key := ari.NewKey(ari.MailboxKey, "mbox1")
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testMailboxCommand(t, m, "Delete", key, nil, cl.Mailbox().Delete)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testMailboxCommand(t, m, "Delete", key, errors.New("err"), cl.Mailbox().Delete)
	})
}

func TestMailboxUpdate(t *testing.T, s Server) {
	key := ari.NewKey(ari.MailboxKey, "mbox1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected error

		m.Mailbox.On("Update", key, 1, 1).Return(expected)

		err := cl.Mailbox().Update(key, 1, 1)

		failed := false
		failed = failed || err == nil && expected != nil
		failed = failed || err != nil && expected == nil
		failed = failed || err != nil && expected != nil && err.Error() != expected.Error() // nolint
		if failed {
			t.Errorf("Expected mailbox %s(%s) to return '%v', got '%v'",
				"Update", "mbox1", expected, err,
			)
		}

		m.Mailbox.AssertCalled(t, "Update", key, 1, 1)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("error")

		m.Mailbox.On("Update", key, 1, 1).Return(expected)

		err := cl.Mailbox().Update(key, 1, 1)

		failed := false
		failed = failed || err == nil && expected != nil
		failed = failed || err != nil && expected == nil
		failed = failed || err != nil && expected != nil && err.Error() != expected.Error()
		if failed {
			t.Errorf("Expected mailbox %s(%s) to return '%v', got '%v'",
				"Update", "mbox1", expected, err,
			)
		}

		m.Mailbox.AssertCalled(t, "Update", key, 1, 1)
	})
}

func TestMailboxData(t *testing.T, s Server) {
	key := ari.NewKey(ari.MailboxKey, "mbox1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected ari.MailboxData
		expected.Name = "mbox1"
		expected.NewMessages = 2
		expected.OldMessages = 3

		m.Mailbox.On("Data", key).Return(&expected, nil)

		data, err := cl.Mailbox().Data(key)
		if err != nil {
			t.Errorf("Unexpected error in remote mailbox Data: %s", err)
		}
		if data == nil {
			t.Errorf("Expected non-nil mailbox data")
		} else {
			failed := data.Name != expected.Name
			failed = failed || data.NewMessages != expected.NewMessages
			failed = failed || data.OldMessages != expected.OldMessages
			if failed {
				t.Errorf("Expected data '%v', got '%v'", expected, data)
			}
		}

		m.Mailbox.AssertCalled(t, "Data", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("error")

		m.Mailbox.On("Data", key).Return(nil, expected)

		_, err := cl.Mailbox().Data(key)

		failed := false
		failed = failed || err == nil && expected != nil
		failed = failed || err != nil && expected == nil
		failed = failed || err != nil && expected != nil && err.Error() != expected.Error()
		if failed {
			t.Errorf("Expected mailbox %s(%s) to return '%v', got '%v'",
				"Data", "mbox1", expected, err,
			)
		}

		m.Mailbox.AssertCalled(t, "Data", key)
	})
}
