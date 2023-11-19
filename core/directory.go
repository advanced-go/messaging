package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// EntryDirectory - collection of Entry
type EntryDirectory interface {
	Add(uri string, c chan Message)
	Count() int
	Uri() []string
	Send(msg Message) error
	Empty()

	get(uri string) *entry
	shutdown()
	empty()
}

// Entry - and entry in an entryDirectory
type entry struct {
	uri string
	c   chan Message
}

type entryDirectory struct {
	m  map[string]*entry
	mu sync.RWMutex
}

// NewEntryDirectory - create a new directory
func NewEntryDirectory() EntryDirectory {
	e := new(entryDirectory)
	e.m = make(map[string]*entry)
	return e
}

func (d *entryDirectory) Add(uri string, c chan Message) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m[uri] = &entry{
		uri: uri,
		c:   c,
	}
}

func (d *entryDirectory) Count() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.m)
}

func (d *entryDirectory) Uri() []string {
	var uri []string
	d.mu.RLock()
	defer d.mu.RUnlock()
	for key, _ := range d.m {
		uri = append(uri, key)
	}
	sort.Strings(uri)
	return uri
}

func (d *entryDirectory) Send(msg Message) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if e, ok := d.m[msg.To]; ok {
		if e.c == nil {
			return errors.New(fmt.Sprintf("entry channel is nil: [%v]", msg.To))
		}
		e.c <- msg
		return nil
	}
	return errors.New(fmt.Sprintf("entry not found: [%v]", msg.To))
}

func (d *entryDirectory) Empty() {
}

func (d *entryDirectory) get(uri string) *entry {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.m[uri]
}

func (d *entryDirectory) shutdown() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, e := range d.m {
		if e.c != nil {
			e.c <- Message{To: e.uri, Event: ShutdownEvent}
		}
	}
}

func (d *entryDirectory) empty() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for key, e := range d.m {
		if e.c != nil {
			close(e.c)
		}
		delete(d.m, key)
	}
}
