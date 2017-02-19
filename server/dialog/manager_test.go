package dialog

import "testing"

func TestMemBind(t *testing.T) {
	m := NewMemManager().(*memManager)

	m.Bind("testDialog", "testType", "testID")

	if len(m.bindings) != 1 {
		t.Errorf("Binding failed; count %d != 1", len(m.bindings))
	}
}

func TestMemBindMultiple(t *testing.T) {
	m := NewMemManager().(*memManager)

	m.Bind("testDialog", "testType", "testID")
	m.Bind("testDialog2", "testType", "testID")

	m.Bind("testDialog", "testType", "testID2")

	m.Bind("testDialog", "testType2", "testID")
	m.Bind("testDialog3", "testType2", "testID")

	if len(m.bindings) != 3 {
		t.Errorf("Binding failed; count %d != 3", len(m.bindings))
	}
}

func TestMemBindDuplicate(t *testing.T) {
	m := NewMemManager().(*memManager)

	m.Bind("testDialog", "testType", "testID")
	m.Bind("testDialog", "testType", "testID")

	if len(m.bindings) != 1 {
		t.Errorf("Binding failed; count %d != 1", len(m.bindings))
	}
}

func TestMemList(t *testing.T) {
	m := NewMemManager().(*memManager)
	m.Bind("testDialog", "testType", "testID")
	m.Bind("testDialog", "testType", "testID2")
	m.Bind("testDialog", "testType2", "testID")
	m.Bind("testDialog2", "testType", "testID")
	m.Bind("testDialog3", "testType2", "testID")

	list := m.List("testType", "testID")

	if len(list) != 2 {
		t.Errorf("Incorrect count %d != 1", len(list))
	}

	var testFound int
	var test2Found int
	for _, d := range list {
		if d == "testDialog" {
			testFound++
		}
		if d == "testDialog2" {
			test2Found++
		}
	}
	if testFound != 1 {
		t.Errorf("Incorrect number of testDialog dialogs: %d != 1", testFound)
	}
	if test2Found != 1 {
		t.Errorf("Incorrect number of testDialog2 dialogs: %d != 1", test2Found)
	}
}
