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
	SendCtrl(msg core.Message) runtime.Status
	SendData(msg core.Message) runtime.Status
}

type directory struct {
	m *sync.Map //[string]*Mailbox
	//mu sync.RWMutex
}

// NewDirectory - create a new directory
func NewDirectory() Directory {
	e := new(directory)
	e.m = new(sync.Map)
	return e
}

func (d *directory) Count() int {
	count := 0
	d.m.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

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

func (d *directory) SendCtrl(msg core.Message) runtime.Status {
	mbox, status := d.get(msg.To)
	if !status.OK() {
		return status.AddLocation(dirSendCtrlLocation)
	}
	if mbox.ctrl == nil {
		return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendCtrlLocation, errors.New(fmt.Sprintf("entry control channel is nil: [%v]", msg.To)))
	}
	mbox.ctrl <- msg
	return runtime.StatusOK()
}

func (d *directory) SendData(msg core.Message) runtime.Status {
	mbox, status := d.get(msg.To)
	if !status.OK() {
		return status.AddLocation(dirSendDataLocation)
	}
	if mbox.data == nil {
		return runtime.NewStatusError(runtime.StatusInvalidContent, dirSendDataLocation, errors.New(fmt.Sprintf("entry data channel is nil: [%v]", msg.To)))
	}
	mbox.data <- msg
	return runtime.StatusOK()
}

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
	if m.ctrl != nil {
		close(m.ctrl)
	}
	if m.data != nil {
		close(m.data)
	}
	d.m.Delete(uri)
	return runtime.StatusOK()
}

func (d *directory) add(m *Mailbox) runtime.Status {
	if m == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, dirAddLocation, errors.New(fmt.Sprintf("invalid argument: mailbox is nil")))
	}
	d.m.Store(m.uri, m)
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
