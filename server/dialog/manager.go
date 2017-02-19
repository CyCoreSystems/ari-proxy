package dialog

import "sync"

// Manager is a dialog manager, which tracks associations between dialogs and entities
type Manager interface {
	// List returns a list of dialogs for the given entity type-ID pair
	List(eType, id string) []string

	// Bind binds the given dialog to an entity type-ID pair
	Bind(dialog, eType, id string)

	// Unbind removes bindings for the given entity type-ID pair
	Unbind(eType, id string)

	// UnbindDialog removes all bindings for the given dialog
	UnbindDialog(dialog string)
}

// Binding is a binding of a Dialog to an entity-type pair
type Binding struct {
	// Dialog is the dialog ID
	Dialog string

	// Type is the entity type
	Type string

	// ID is the unique identifier of the entity
	ID string
}

func bindingHash(eType, id string) string {
	return eType + ":" + id
}

type memManager struct {
	bindings map[string][]string

	mu sync.RWMutex
}

// NewMemManager returns a new in-memory dialog manager
func NewMemManager() Manager {
	return &memManager{
		bindings: make(map[string][]string),
	}
}

func (m *memManager) List(eType, id string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list, ok := m.bindings[bindingHash(eType, id)]
	if !ok {
		return nil
	}
	return list
}

func (m *memManager) Bind(dialog, eType, id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	list, ok := m.bindings[bindingHash(eType, id)]
	if !ok {
		list = []string{dialog}
	} else {
		// Don't add the binding if it is already there
		for _, b := range list {
			if b == dialog {
				return
			}
		}

		// Add the dialog
		list = append(list, dialog)
	}
	m.bindings[bindingHash(eType, id)] = list
}

func (m *memManager) Unbind(eType, id string) {
	m.mu.Lock()
	delete(m.bindings, bindingHash(eType, id))
	m.mu.Unlock()
}

func (m *memManager) UnbindDialog(dialog string) {

	var dialogIndex int

	m.mu.Lock()
	for _, v := range m.bindings {
		dialogIndex = -1
		for i, d := range v {
			if d == dialog {
				dialogIndex = i
			}
		}
		if dialogIndex >= 0 {
			v[dialogIndex] = v[len(v)-1]
			v = v[:len(v)-1]
		}
	}
	m.mu.Unlock()

	m.mu.Unlock()
}
