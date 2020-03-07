package cluster

import (
	"strings"
	"testing"
	"time"

	rid "github.com/CyCoreSystems/ari-rid"
)

func TestHash(t *testing.T) {
	id := rid.New("")
	app := rid.New("")

	testID, testApp := dehash(hash(id, app))
	if id != testID {
		t.Error("Asterisk IDs do not match")
	}
	if app != testApp {
		t.Error("ARI Applications do not match")
	}
}

func TestAll(t *testing.T) {
	c := New()
	c.Update("A1", "TestApp")
	c.Update("A2", "TestApp")

	list := c.All(0)
	if len(list) != 2 {
		t.Errorf("Incorrect number of cluster members: %d != 2", len(list))
	}

	c.Update("B1", "TestApp2")
	c.Update("B2", "TestApp2")

	list = c.All(0)
	if len(list) != 4 {
		t.Errorf("Incorrect number of cluster members: %d != 4", len(list))
	}
}

func TestApp(t *testing.T) {
	c := New()
	c.Update("A1", "TestApp")
	c.Update("A2", "TestApp")
	c.Update("B1", "TestApp2")
	c.Update("B2", "TestApp2")

	list := c.App("TestApp", 0)
	if len(list) != 2 {
		t.Errorf("Incorrect number of cluster members: %d != 2", len(list))
	}
	if list[1].App != "TestApp" {
		t.Errorf("Incorrect app: %s != TestApp", list[0].App)
	}
	if !strings.HasPrefix(list[0].ID, "A") {
		t.Errorf("Incorrect ID: %s does not begin with A", list[1].ID)
	}

	list = c.App("TestApp2", 0)
	if len(list) != 2 {
		t.Errorf("Incorrect number of cluster members: %d != 2", len(list))
	}
	if list[0].App != "TestApp2" {
		t.Errorf("Incorrect app: %s != TestApp2", list[0].App)
	}
	if !strings.HasPrefix(list[1].ID, "B") {
		t.Errorf("Incorrect ID: %s does not begin with B", list[1].ID)
	}
}

func TestPurge(t *testing.T) {
	c := New()
	c.Update("A1", "TestApp")
	c.Update("A2", "TestApp")
	c.Update("B1", "TestApp2")
	c.Update("B2", "TestApp2")

	list := c.All(0)
	if len(list) != 4 {
		t.Errorf("Incorrect number of cluster members: %d != 4", len(list))
	}

	c.Purge(50 * time.Millisecond)

	list = c.All(0)
	if len(list) != 4 {
		t.Errorf("Incorrect number of cluster members: %d != 4", len(list))
	}

	time.Sleep(50 * time.Millisecond)

	c.Update("B1", "TestApp2")

	c.Purge(45 * time.Millisecond)

	list = c.All(0)
	if len(list) != 1 {
		t.Errorf("Incorrect number of cluster members: %d != 1", len(list))
	}

	c.Update("A2", "TestApp")

	list = c.All(0)
	if len(list) != 2 {
		t.Errorf("Incorrect number of cluster members: %d != 2", len(list))
	}
}
