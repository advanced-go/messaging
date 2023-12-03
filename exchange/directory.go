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
	dirSendCtrlLocation = PkgPath + ":Directory/SendCtrl"
	dirSendDataLocation = PkgPath + ":Directory/SendData"
)

// Directory - exchange directory
type Directory interface {
	Count() int
	List() []string
	SendCtrl(msg core.Message) runtime.Status
	SendData(msg core.Message) runtime.Status
	Shutdown()
}

type directory struct {
	m *sync.Map //[string]*Mailbox
	//mu sync.RWMutex
}

// NewDirectory - create a new directory
func NewDirectory() Directory {
	e := new(directory)
	e.m = new(sync.Map) //make(map[string]*Mailbox)
	return e
}

func (d *directory) Count() int {
	//d.mu.RLock()
	//defer d.mu.RUnlock()
	count := 0
	d.m.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

func (d *directory) List() []string {
	var uri []string
	//d.mu.RLock()
	//defer d.mu.RUnlock()
	d.m.Range(func(key, value any) bool {
		if str, ok := key.(string); ok {
			uri = append(uri, str)
		}
		return true
	})
	//for key, _ := range d.m {
	//}
	sort.Strings(uri)
	return uri
}

func (d *directory) SendCtrl(msg core.Message) runtime.Status {
	v, ok := d.m.Load(msg.To)
	if !ok {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, dirSendCtrlLocation, errors.New(fmt.Sprintf("entry not found: [%v]", msg.To)))
	}
	mbox, ok2 := v.(*Mailbox)
	if !ok2 {
		return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendCtrlLocation, errors.New(fmt.Sprintf("invalid Mailbox type: [%v]", msg.To)))
	}
	if mbox.ctrl == nil {
		return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendCtrlLocation, errors.New(fmt.Sprintf("entry control channel is nil: [%v]", msg.To)))
	}
	mbox.ctrl <- msg
	return runtime.StatusOK()
	//d.mu.RLock()
	//defer d.mu.RUnlock()
	/*
		if e, ok := d.m[msg.To]; ok {
			if e.ctrl == nil {
				return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendLocation, errors.New(fmt.Sprintf("entry command channel is nil: [%v]", msg.To)))
			}
			e.ctrl <- msg
			return runtime.StatusOK()
		}
		return runtime.NewStatusError(runtime.StatusInvalidArgument, dirSendLocation, errors.New(fmt.Sprintf("entry not found: [%v]", msg.To)))
		return runtime.StatusOK()
	*/
}

func (d *directory) SendData(msg core.Message) runtime.Status {
	v, ok := d.m.Load(msg.To)
	if !ok {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, dirSendDataLocation, errors.New(fmt.Sprintf("entry not found: [%v]", msg.To)))
	}
	mbox, ok2 := v.(*Mailbox)
	if !ok2 {
		return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendDataLocation, errors.New(fmt.Sprintf("invalid Mailbox type: [%v]", msg.To)))
	}
	if mbox.data == nil {
		return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendDataLocation, errors.New(fmt.Sprintf("entry data channel is nil: [%v]", msg.To)))
	}
	mbox.data <- msg
	return runtime.StatusOK()

	/*
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
	*/

}

func (d *directory) Shutdown() {
	//d.mu.RLock()
	//defer d.mu.RUnlock()
	//for _, e := range d.m {
	//	if e.ctrl != nil {
	//		e.ctrl <- core.Message{To: e.uri, Event: core.ShutdownEvent}
	//	}
	//}
}

func add(dir Directory, m *Mailbox) error {
	if dir == nil {
		return errors.New(fmt.Sprintf("invalid argument: directory is nil"))
	}
	if m == nil {
		return errors.New(fmt.Sprintf("invalid argument: mailbox is nil"))
	}
	if d, ok := any(dir).(*directory); ok {
		//d.mu.Lock()
		//defer d.mu.Unlock()
		d.m.Store(m.uri, m)
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
	d, ok := any(dir).(*directory)
	if !ok {
		return nil, errors.New("invalid argument: Directory is not of type *directory")
	}
	//d.mu.Lock()
	//defer d.mu.Unlock()
	v, ok1 := d.m.Load(uri)
	if !ok1 {
		return nil, errors.New("invalid URI: directory entry not found")
	}
	if mbox, ok2 := v.(*Mailbox); ok2 {
		return mbox, nil
	}
	return nil, errors.New("invalid Mailbox type")
}
