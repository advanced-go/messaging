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
	dirSendLocation = PkgPath + ":Directory/Send"
)

// Directory - exchange directory
type Directory interface {
	Count() int
	List() []string
	SendCmd(msg core.Message) runtime.Status
	SendData(msg core.Message) runtime.Status
	Shutdown()
}

type directory struct {
	m  map[string]*Mailbox
	mu sync.RWMutex
}

// NewDirectory - create a new directory
func NewDirectory() Directory {
	e := new(directory)
	e.m = make(map[string]*Mailbox)
	return e
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

func (d *directory) SendCmd(msg core.Message) runtime.Status {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if e, ok := d.m[msg.To]; ok {
		if e.cmd == nil {
			return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendLocation, errors.New(fmt.Sprintf("entry command channel is nil: [%v]", msg.To)))
		}
		e.cmd <- msg
		return runtime.StatusOK()
	}
	return runtime.NewStatusError(runtime.StatusInvalidArgument, dirSendLocation, errors.New(fmt.Sprintf("entry not found: [%v]", msg.To)))
}

func (d *directory) SendData(msg core.Message) runtime.Status {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if e, ok := d.m[msg.To]; ok {
		if e.data == nil {
			return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendLocation, errors.New(fmt.Sprintf("entry data channel is nil: [%v]", msg.To)))
		}
		e.data <- msg
		return runtime.StatusOK()
	}
	return runtime.NewStatusError(runtime.StatusInvalidArgument, dirSendLocation, errors.New(fmt.Sprintf("entry not found: [%v]", msg.To)))
}

func (d *directory) Shutdown() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, e := range d.m {
		if e.cmd != nil {
			e.cmd <- core.Message{To: e.uri, Event: core.ShutdownEvent}
		}
	}
}

func add(dir Directory, m *Mailbox) error {
	if dir == nil {
		return errors.New(fmt.Sprintf("invalid argument: directory is nil"))
	}
	if m == nil {
		return errors.New(fmt.Sprintf("invalid argument: mailbox is nil"))
	}
	if d, ok := any(dir).(*directory); ok {
		d.mu.Lock()
		defer d.mu.Unlock()
		d.m[m.uri] = m
	}
	return nil
}

func get(dir Directory, uri string) (*Mailbox, error) {
	if dir == nil {
		return nil, errors.New("invalid argument: directory is nil")
	}
	if len(uri) == 0 {
		return nil, errors.New("invalid argument: uri is empty")
	}
	if d, ok := any(dir).(*directory); ok {
		d.mu.Lock()
		defer d.mu.Unlock()
		return d.m[uri], nil
	}
	return nil, errors.New("invalid argument: Directory is not of type *directory")
}
