package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
)

func TestMailboxList(t *testing.T, s Server) {
	runTest("empty", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Mailbox.On("List").Return([]ari.MailboxHandle{}, nil)

		ret, err := cl.Mailbox().List()
		if err != nil {
			t.Errorf("Unexpected error in remote List call")
		}
		if len(ret) != 0 {
			t.Errorf("Expected return length to be 0, got %d", len(ret))
		}

		m.Shutdown()

		m.Mailbox.AssertCalled(t, "List")
	})

	runTest("nonEmpty", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var h1 = &mocks.MailboxHandle{}
		var h2 = &mocks.MailboxHandle{}

		h1.On("ID").Return("h1")
		h2.On("ID").Return("h2")

		m.Mailbox.On("List").Return([]ari.MailboxHandle{h1, h2}, nil)

		ret, err := cl.Mailbox().List()
		if err != nil {
			t.Errorf("Unexpected error in remote List call")
		}
		if len(ret) != 2 {
			t.Errorf("Expected return length to be 2, got %d", len(ret))
		}

		m.Shutdown()

		m.Mailbox.AssertCalled(t, "List")
		h1.AssertCalled(t, "ID")
		h2.AssertCalled(t, "ID")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Mailbox.On("List").Return(nil, errors.New("unknown error"))

		ret, err := cl.Mailbox().List()
		if err == nil {
			t.Errorf("Expected error in remote List call")
		}
		if len(ret) != 0 {
			t.Errorf("Expected return length to be 0, got %d", len(ret))
		}

		m.Shutdown()

		m.Mailbox.AssertCalled(t, "List")
	})
}

func testMailboxCommand(t *testing.T, m *mock, name string, id string, expected error, fn func(string) error) {
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
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testMailboxCommand(t, m, "Delete", "mbox1", nil, cl.Mailbox().Delete)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testMailboxCommand(t, m, "Delete", "mbox1", errors.New("err"), cl.Mailbox().Delete)
	})
}

func TestMailboxUpdate(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected error

		m.Mailbox.On("Update", "mbox1", 1, 1).Return(expected)

		err := cl.Mailbox().Update("mbox1", 1, 1)

		failed := false
		failed = failed || err == nil && expected != nil
		failed = failed || err != nil && expected == nil
		failed = failed || err != nil && expected != nil && err.Error() != expected.Error()
		if failed {
			t.Errorf("Expected mailbox %s(%s) to return '%v', got '%v'",
				"Update", "mbox1", expected, err,
			)
		}

		m.Mailbox.AssertCalled(t, "Update", "mbox1", 1, 1)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = errors.New("error")

		m.Mailbox.On("Update", "mbox1", 1, 1).Return(expected)

		err := cl.Mailbox().Update("mbox1", 1, 1)

		failed := false
		failed = failed || err == nil && expected != nil
		failed = failed || err != nil && expected == nil
		failed = failed || err != nil && expected != nil && err.Error() != expected.Error()
		if failed {
			t.Errorf("Expected mailbox %s(%s) to return '%v', got '%v'",
				"Update", "mbox1", expected, err,
			)
		}

		m.Mailbox.AssertCalled(t, "Update", "mbox1", 1, 1)
	})
}

func TestMailboxData(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var expected ari.MailboxData
		expected.Name = "mbox1"
		expected.NewMessages = 2
		expected.OldMessages = 3

		m.Mailbox.On("Data", "mbox1").Return(&expected, nil)

		data, err := cl.Mailbox().Data("mbox1")
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

		m.Mailbox.AssertCalled(t, "Data", "mbox1")
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = errors.New("error")

		m.Mailbox.On("Data", "mbox1").Return(nil, expected)

		_, err := cl.Mailbox().Data("mbox1")

		failed := false
		failed = failed || err == nil && expected != nil
		failed = failed || err != nil && expected == nil
		failed = failed || err != nil && expected != nil && err.Error() != expected.Error()
		if failed {
			t.Errorf("Expected mailbox %s(%s) to return '%v', got '%v'",
				"Data", "mbox1", expected, err,
			)
		}

		m.Mailbox.AssertCalled(t, "Data", "mbox1")
	})
}
