package plugin

import (
	"fmt"
	"sync"
)

type manager struct {
	sync.Mutex
	plugins    []Plugin
	registered map[string]bool
}

// global plugin manager
var defaultManager = newManager()

func newManager() *manager {
	return &manager{
		registered: make(map[string]bool),
	}
}

func (m *manager) Plugins() []Plugin {
	m.Lock()
	defer m.Unlock()
	return m.plugins
}

func (m *manager) Register(plugin Plugin) error {
	m.Lock()
	defer m.Unlock()

	name := plugin.String()

	if m.registered[name] {
		return fmt.Errorf("plugin with name %s already registered", name)
	}

	m.registered[name] = true
	m.plugins = append(m.plugins, plugin)
	return nil
}

func (m *manager) isRegistered(plugin Plugin) bool {
	m.Lock()
	defer m.Unlock()

	return m.registered[plugin.String()]
}
