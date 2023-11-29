package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"sort"
	"sync"
)

const (
	dirSendLocation = PkgPath + "/Directory/Send"
)

// Directory - exchange directory
type Directory interface {
	Add(uri string, c chan core.Message)
	Count() int
	List() []string
	Send(msg core.Message) runtime.Status
	Shutdown()
}

type entry struct {
	uri string
	c   chan core.Message
}

type directory struct {
	m  map[string]*entry
	mu sync.RWMutex
}

// NewDirectory - create a new directory
func NewDirectory() Directory {
	e := new(directory)
	e.m = make(map[string]*entry)
	return e
}

func (d *directory) Add(uri string, c chan core.Message) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m[uri] = &entry{
		uri: uri,
		c:   c,
	}
}

func (d *directory) Count() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.m)
}

func (d *directory) List() []string {
	var uri []string
	d.mu.RLock()
	defer d.mu.RUnlock()
	for key, _ := range d.m {
		uri = append(uri, key)
	}
	sort.Strings(uri)
	return uri
}

func (d *directory) Send(msg core.Message) runtime.Status {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if e, ok := d.m[msg.To]; ok {
		if e.c == nil {
			return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendLocation, errors.New(fmt.Sprintf("entry channel is nil: [%v]", msg.To)))
		}
		e.c <- msg
		return runtime.StatusOK()
	}
	return runtime.NewStatusError(runtime.StatusInvalidArgument, dirSendLocation, errors.New(fmt.Sprintf("entry not found: [%v]", msg.To)))
}

func (d *directory) Shutdown() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, e := range d.m {
		if e.c != nil {
			e.c <- core.Message{To: e.uri, Event: core.ShutdownEvent}
		}
	}
}
