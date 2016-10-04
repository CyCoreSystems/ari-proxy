package ariproxy

import "sync"

// index of instances used for routing events
type instanceCache struct {
	cache map[string]*Instance
	lock  sync.RWMutex
}

func (ic *instanceCache) Add(id string, i *Instance) {

	i.Dialog.Objects.Add(id)
	ic.lock.Lock()
	defer ic.lock.Unlock()

	ic.cache[id] = i
}

func (ic *instanceCache) Remove(id string, i *Instance) {
	i.Dialog.Objects.Remove(id)

	ic.lock.Lock()
	defer ic.lock.Unlock()

	delete(ic.cache, id)
}

func (ic *instanceCache) Find(id string) *Instance {
	ic.lock.RLock()
	defer ic.lock.RUnlock()

	i, _ := ic.cache[id]
	return i
}

func (ic *instanceCache) RemoveAll(i *Instance) {
	i.Dialog.Objects.Clear()

	ic.lock.Lock()
	defer ic.lock.Unlock()

	for id, in := range ic.cache {
		if in.Dialog.ID == i.Dialog.ID {
			delete(ic.cache, id)
		}
	}

}
