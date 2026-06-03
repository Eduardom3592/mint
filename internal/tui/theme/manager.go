package theme

import (
	"database/sql"
	"fmt"
	"sync"
)

type Manager struct {
	mu      sync.RWMutex
	current *Theme
	themes  map[string]Theme
	order   []string
	db      *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	m := &Manager{
		themes: make(map[string]Theme),
		order:  make([]string, 0),
		db:     db,
	}

	for _, t := range BuiltInThemes() {
		m.register(t)
	}

	if saved, err := m.loadSavedTheme(); err == nil && saved != nil {
		m.current = saved
	} else {
		first, _ := m.Get("Mint")
		m.current = first
	}

	return m
}

func (m *Manager) register(t Theme) {
	key := t.Name
	if _, exists := m.themes[key]; !exists {
		m.themes[key] = t
		m.order = append(m.order, key)
	}
}

func (m *Manager) List() []Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]Theme, 0, len(m.order))
	for _, name := range m.order {
		list = append(list, m.themes[name])
	}
	return list
}

func (m *Manager) Get(name string) (*Theme, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, ok := m.themes[name]
	if !ok {
		return nil, fmt.Errorf("theme %q not found", name)
	}
	return &t, nil
}

func (m *Manager) Current() *Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

func (m *Manager) Set(name string) (*Theme, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.themes[name]
	if !ok {
		return nil, fmt.Errorf("theme %q not found", name)
	}
	m.current = &t

	m.persist(name)

	return &t, nil
}

func (m *Manager) RegisterCustom(t Theme) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.themes[t.Name] = t
	m.order = append(m.order, t.Name)
}

func (m *Manager) persist(name string) {
	if m.db == nil {
		return
	}
	_, _ = m.db.Exec(`INSERT INTO settings (key, value) VALUES ('theme', ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, name)
}

func (m *Manager) loadSavedTheme() (*Theme, error) {
	if m.db == nil {
		return nil, nil
	}

	row := m.db.QueryRow(`SELECT value FROM settings WHERE key = 'theme'`)
	var name string
	if err := row.Scan(&name); err != nil {
		return nil, err
	}

	t, ok := m.themes[name]
	if !ok {
		return nil, fmt.Errorf("saved theme %q not available", name)
	}

	return &t, nil
}
