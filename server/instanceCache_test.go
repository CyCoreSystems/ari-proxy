package ariproxy

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/session"
)

func TestInstanceCacheInit(t *testing.T) {
	var ic instanceCache
	ic.Init()

	if ic.cache == nil {
		t.Errorf("ic.Init did not create cache")
	}

	if ic.associations == nil {
		t.Error("ic.Init did not create associations")
	}
}

func TestInstanceCacheAdd(t *testing.T) {
	var ic instanceCache
	ic.Init()

	ic.Add("object1", nil)

	if x := len(ic.cache); x != 0 {
		t.Errorf("ic.Add('object1', nil), cache should be empty, is %d", x)
	}

	if x := len(ic.associations); x != 0 {
		t.Errorf("ic.Add('object1', nil), associations should be empty, is %d", x)
	}

	var i Instance
	i.Dialog = session.NewDialog("object1", nil)

	ic.Add("object1", &i)

	if x := len(ic.cache); x != 1 {
		t.Errorf("ic.Add('object1', i), cache should be 1, is %d", x)
	}

	if x := len(ic.associations); x != 1 {
		t.Errorf("ic.Add('object1', i), associations should be 1, is %d", x)
	}
}

func TestInstanceCacheFind(t *testing.T) {
	var ic instanceCache
	ic.Init()

	ix := ic.Find("object1")

	if x := len(ix); x != 0 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			0, len(ix))
	}

	var i1 Instance
	i1.Dialog = session.NewDialog("dialog1", nil)

	ic.Add("object1", &i1)

	ix = ic.Find("object1")

	if x := len(ix); x != 1 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			1, len(ix))
	}

	var i2 Instance
	i2.Dialog = session.NewDialog("dialog2", nil)

	ic.Add("object1", &i2)

	ix = ic.Find("object1")

	if x := len(ix); x != 2 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			2, len(ix))
	}

	ic.Add("object2", &i2)

	ix = ic.Find("object2")

	if x := len(ix); x != 1 {
		t.Errorf("ic.Find('object2') => expected len %d, was len %d",
			2, len(ix))
	}
}

func TestInstanceCacheRemoveAll(t *testing.T) {
	var ic instanceCache
	ic.Init()
	ic.RemoveAll(nil)

	var i1 Instance
	i1.Dialog = session.NewDialog("dialog1", nil)

	ic.Add("object1", &i1)

	ix := ic.Find("object1")

	if x := len(ix); x != 1 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			1, len(ix))
	}

	ic.RemoveAll(&i1)

	ix = ic.Find("object1")

	if x := len(ix); x != 0 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			0, len(ix))
	}

}

func TestInstanceCacheRemoveObject(t *testing.T) {
	var ic instanceCache
	ic.Init()

	ix := ic.Find("object1")

	if x := len(ix); x != 0 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			0, len(ix))
	}

	var i1 Instance
	i1.Dialog = session.NewDialog("dialog1", nil)

	ic.Add("object1", &i1)

	ix = ic.Find("object1")

	if x := len(ix); x != 1 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			1, len(ix))
	}

	var i2 Instance
	i2.Dialog = session.NewDialog("dialog2", nil)

	ic.Add("object1", &i2)

	ix = ic.Find("object1")
	if x := len(ix); x != 2 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			2, len(ix))
	}

	ic.Add("object2", &i2)

	ix = ic.Find("object2")
	if x := len(ix); x != 1 {
		t.Errorf("ic.Find('object2') => expected len %d, was len %d",
			2, len(ix))
	}

	ic.RemoveObject("object2", nil)
	ic.RemoveObject("object2", &i2)
	ic.RemoveObject("object2", &i1)

	if x := len(ix); x != 1 {
		t.Errorf("ic.Find('object2') => expected len %d, was len %d",
			2, len(ix))
	}

	ix = ic.Find("object1")
	if x := len(ix); x != 2 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			2, len(ix))
	}

	ix = ic.Find("object2")
	if x := len(ix); x != 0 {
		t.Errorf("ic.Find('object2') => expected len %d, was len %d",
			0, len(ix))
	}

	ic.RemoveObject("object1", nil)
	ic.RemoveObject("object1", &i2)
	ic.RemoveObject("object1", &i1)

	ix = ic.Find("object1")
	if x := len(ix); x != 0 {
		t.Errorf("ic.Find('object1') => expected len %d, was len %d",
			0, len(ix))
	}

	ix = ic.Find("object2")
	if x := len(ix); x != 0 {
		t.Errorf("ic.Find('object2') => expected len %d, was len %d",
			0, len(ix))
	}

}
