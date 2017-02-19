package ariproxy

import (
	"sort"
	"sync"
)

// index of instances used for routing events
type instanceCache struct {
	cache        map[string]*Instance // instance container
	associations map[string][]string  // object ID to instance ID associations

	lock sync.RWMutex
}

func (ic *instanceCache) Init() {
	ic.cache = make(map[string]*Instance)
	ic.associations = make(map[string][]string)
}

func (ic *instanceCache) Add(objectID string, i *Instance) {

	if i == nil {
		return
	}

	i.Dialog.Objects.Add(objectID)

	ic.lock.Lock()
	defer ic.lock.Unlock()

	dialogID := i.Dialog.ID

	// store the instance
	ic.cache[dialogID] = i

	// associate the dialog ID with the object ID
	if _, ok := ic.associations[objectID]; !ok {
		ic.associations[objectID] = make([]string, 0)
	}
	ic.associations[objectID] = append(ic.associations[objectID], dialogID)
	sort.Strings(ic.associations[objectID])
}

func (ic *instanceCache) RemoveObject(objectID string, i *Instance) {
	if i == nil {
		return
	}

	i.Dialog.Objects.Remove(objectID)

	ic.lock.Lock()
	defer ic.lock.Unlock()

	ic.removeObject(objectID, i)
}

func (ic *instanceCache) removeObject(objectID string, i *Instance) {
	if i == nil {
		return
	}

	dialogID := i.Dialog.ID

	if lx, ok := ic.associations[objectID]; ok && len(lx) != 0 {
		idx := sort.SearchStrings(lx, dialogID)
		if lx[idx] == dialogID {
			lx = append(lx[:idx], lx[idx+1:]...)
		}
		ic.associations[objectID] = lx
	}
}

func (ic *instanceCache) Find(objectID string) (ix []*Instance) {
	ic.lock.RLock()
	defer ic.lock.RUnlock()

	dialogs, _ := ic.associations[objectID]
	for _, dialog := range dialogs {
		ix = append(ix, ic.cache[dialog])
	}

	return
}

func (ic *instanceCache) RemoveAll(i *Instance) {
	if i == nil {
		return
	}

	objectIDs := i.Dialog.Objects.Items()
	i.Dialog.Objects.Clear()

	ic.lock.Lock()
	defer ic.lock.Unlock()

	dialogID := i.Dialog.ID
	delete(ic.cache, dialogID)

	for _, objectID := range objectIDs {
		ic.removeObject(objectID, i)
	}
}
