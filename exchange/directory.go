package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"net/http"
	"sort"
	"sync"
)

const (
	dirSendCtrlLocation = PkgPath + ":Directory/SendCtrl"
	dirSendDataLocation = PkgPath + ":Directory/SendData"
	dirGetLocation      = PkgPath + ":Directory/get"
	dirAddLocation      = PkgPath + ":Directory/add"
)

// Directory - exchange directory
type Directory interface {
	Count() int
	List() []string
	Add(m *Mailbox) runtime.Status
	SendCtrl(msg core.Message) runtime.Status
	SendData(msg core.Message) runtime.Status
}

type directory struct {
	m *sync.Map
}

// NewDirectory - create a new directory
func NewDirectory() Directory {
	e := new(directory)
	e.m = new(sync.Map)
	return e
}

// Count - number of items in the sync map
func (d *directory) Count() int {
	count := 0
	d.m.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

// List - a list of item uri's
func (d *directory) List() []string {
	var uri []string
	d.m.Range(func(key, value any) bool {
		if str, ok := key.(string); ok {
			uri = append(uri, str)
		}
		return true
	})
	sort.Strings(uri)
	return uri
}

// SendCtrl - send a message to the select item's control channel
func (d *directory) SendCtrl(msg core.Message) runtime.Status {
	// TO DO : authenticate shutdown control message
	if msg.Event == core.ShutdownEvent {
		return runtime.StatusOK()
	}
	mbox, status := d.get(msg.To)
	if !status.OK() {
		return status.AddLocation(dirSendCtrlLocation)
	}
	mbox.SendCtrl(msg)
	return runtime.StatusOK()
}

// SendData - send a message to the item's data channel
func (d *directory) SendData(msg core.Message) runtime.Status {
	mbox, status := d.get(msg.To)
	if !status.OK() {
		return status.AddLocation(dirSendDataLocation)
	}
	mbox.SendData(msg)
	return runtime.StatusOK()
}

// Add - add a mailbox
func (d *directory) Add(m *Mailbox) runtime.Status {
	if m == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, dirAddLocation, errors.New("invalid argument: mailbox is nil"))
	}
	if len(m.uri) == 0 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, dirAddLocation, errors.New("invalid argument: mailbox uri is empty"))
	}
	if m.ctrl == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, dirAddLocation, errors.New("invalid argument: mailbox command channel is nil"))
	}
	_, ok := d.m.Load(m.uri)
	if ok {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, dirAddLocation, errors.New(fmt.Sprintf("invalid argument: directory mailbox already exists: [%v]", m.uri)))
	}
	d.m.Store(m.uri, m)
	m.unregister = func() {
		d.m.Delete(m.uri)
	}
	return runtime.StatusOK()
}

func (d *directory) get(uri string) (*Mailbox, runtime.Status) {
	if len(uri) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, dirGetLocation, errors.New("invalid argument: uri is empty"))
	}
	v, ok1 := d.m.Load(uri)
	if !ok1 {
		return nil, runtime.NewStatusError(http.StatusNotFound, dirGetLocation, errors.New(fmt.Sprintf("invalid URI: directory mailbox not found [%v]", uri)))
	}
	if mbox, ok2 := v.(*Mailbox); ok2 {
		return mbox, runtime.StatusOK()
	}
	return nil, runtime.NewStatusError(runtime.StatusInvalidContent, dirGetLocation, errors.New("invalid Mailbox type"))
}

// Shutdown - close an item's mailbox
func (d *directory) Shutdown(msg core.Message) runtime.Status {
	// TO DO: add authentication
	return runtime.StatusOK() //d.shutdown(msg.To)
}

/*
func (d *directory) shutdown(uri string) runtime.Status {
	//d.mu.RLock()
	//defer d.mu.RUnlock()
	//for _, e := range d.m {
	//	if e.ctrl != nil {
	//		e.ctrl <- core.Message{To: e.uri, Event: core.ShutdownEvent}
	//	}
	//}
	m, status := d.get(uri)
	if !status.OK() {
		return status
	}
	if m.data != nil {
		close(m.data)
	}
	if m.ctrl != nil {
		close(m.ctrl)
	}
	d.m.Delete(uri)
	return runtime.StatusOK()
}
*/
