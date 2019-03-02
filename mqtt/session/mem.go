package session

import (
	"fmt"
	"sync"
)

var _ SessionsProvider = (*memProvider)(nil)

type memProvider struct {
	st map[string]*Session
	mu sync.RWMutex
}

// NewMemProvider creates a in-memory sessions storage provider.
func NewMemProvider() *SessionsProvider {
	p := SessionsProvider(&memProvider{
		st: make(map[string]*Session),
	})
	return &p
}

func (p *memProvider) New(id string) (*Session, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if id == "" {
		id = generateSessionId()
	}
	p.st[id] = &Session{id: id}
	return p.st[id], nil
}

func (p *memProvider) Get(id string) (*Session, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	sess, ok := p.st[id]
	if !ok {
		return nil, fmt.Errorf("store/Get: No session found for key %s", id)
	}

	return sess, nil
}

func (p *memProvider) Delete(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.st, id)
}

func (p *memProvider) Save(id string) error {
	return nil
}

func (p *memProvider) Count() int {
	return len(p.st)
}

func (p *memProvider) Close() error {
	p.st = make(map[string]*Session)
	return nil
}
