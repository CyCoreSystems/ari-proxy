package session

import (
	"sort"
	"sync"
)

// Objects tracks a list of object IDs that are associated with the dialog
type Objects struct {
	ids   []string
	mutex sync.RWMutex
}

// Contains finds the id, if it exists
func (o *Objects) Contains(id string) (int, bool) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	idx := sort.SearchStrings(o.ids, id)

	if idx == len(o.ids) || o.ids[idx] != id {
		return -1, false
	}

	return idx, true
}

// Add adds the given object
func (o *Objects) Add(id string) bool {

	if _, ok := o.Contains(id); ok {
		return false
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.ids = append(o.ids, id)
	sort.Strings(o.ids)

	return true
}

// Remove removes the given object, if it exists
func (o *Objects) Remove(id string) bool {
	idx, ok := o.Contains(id)
	if !ok {
		return false
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.ids = append(o.ids[:idx], o.ids[idx+1:]...)

	return true
}

// Clear removes all the objects
func (o *Objects) Clear() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.ids = make([]string, 0)
}

// Items returns the list of items
func (o *Objects) Items() []string {
	return o.ids
}
