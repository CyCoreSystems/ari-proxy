package dialog

import "testing"

func TestMemBind(t *testing.T) {
	m := NewMemManager().(*memManager)

	m.Bind("testDialog", "testType", "testID")

	if len(m.bindings) != 1 {
		t.Errorf("Binding failed; count %d != 1", len(m.bindings))
	}
}

func TestMemUnbindDialog(t *testing.T) {
	m := NewMemManager().(*memManager)

	m.Bind("testDialog", "testType", "testID")

	if len(m.bindings) != 1 {
		t.Errorf("Binding failed; count %d != 1", len(m.bindings))
	}

	m.UnbindDialog("testDialog")

	// The unbinding of a dialog doesn't clear the testType/testID
	if len(m.bindings) != 1 {
		t.Errorf("Unbinding Dialog failed; count %d != 0", len(m.bindings))
	}

	// But it should make the testType/testID empty
	if i := len(m.List("testType", "testID")); i != 0 {
		t.Errorf("List('testType','testID'); count %d != 0", i)
	}
}

func TestMemBindUnbind(t *testing.T) {
	m := NewMemManager().(*memManager)

	m.Bind("testDialog", "testType", "testID")

	if len(m.bindings) != 1 {
		t.Errorf("Binding failed; count %d != 1", len(m.bindings))
	}

	m.Unbind("testType", "testID")

	if len(m.bindings) != 0 {
		t.Errorf("Unbinding failed; count %d != 0", len(m.bindings))
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

func TestMemBindUnbindMultiple(t *testing.T) {
	m := NewMemManager().(*memManager)

	m.Bind("testDialog", "testType", "testID")
	m.Bind("testDialog2", "testType", "testID")

	m.Bind("testDialog", "testType", "testID2")

	m.Bind("testDialog", "testType2", "testID")
	m.Bind("testDialog3", "testType2", "testID")

	if len(m.bindings) != 3 {
		t.Errorf("Binding failed; count %d != 3", len(m.bindings))
	}

	m.Unbind("testType", "testID")

	if len(m.bindings) != 2 {
		t.Errorf("Unbinding failed; count %d != 2", len(m.bindings))
	}

	if len(m.List("testType", "testID")) != 0 {
		t.Errorf("Unbinding failed; List('testType','testID') => len %d != 2", len(m.bindings))
	}

	m.Unbind("testType2", "testID")

	if len(m.bindings) != 1 {
		t.Errorf("Binding failed; count %d != 1", len(m.bindings))
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
