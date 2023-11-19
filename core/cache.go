package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// MessageCache - message cache by uri
type MessageCache interface {
	Count() int
	Filter(event string, code int, include bool) []string
	Include(event string, status int) []string
	Exclude(event string, status int) []string
	Add(msg Message) error
	Get(uri string) (Message, error)
	Uri() []string
}

type messageCache struct {
	m  map[string]Message
	mu sync.RWMutex
}

// NewMessageCache - create a message cache
func NewMessageCache() MessageCache {
	c := new(messageCache)
	c.m = make(map[string]Message)
	return c
}

func (r *messageCache) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, _ = range r.m {
		count++
	}
	return count
}

func (r *messageCache) Filter(event string, code int, include bool) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var uri []string
	for u, resp := range r.m {
		if include {
			if resp.Status != nil && resp.Status.Code() == code && resp.Event == event {
				uri = append(uri, u)
			}
		} else {
			if resp.Status != nil && resp.Status.Code() != code || resp.Event != event {
				uri = append(uri, u)
			}
		}
	}
	sort.Strings(uri)
	return uri
}

func (r *messageCache) Include(event string, status int) []string {
	return r.Filter(event, status, true)
}

func (r *messageCache) Exclude(event string, status int) []string {
	return r.Filter(event, status, false)
}

func (r *messageCache) Add(msg Message) error {
	if msg.From == "" {
		return errors.New("invalid argument: message from is empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[msg.From]; !ok {
		r.m[msg.From] = msg
	}
	return nil
}

func (r *messageCache) Get(uri string) (Message, error) {
	if uri == "" {
		return Message{}, errors.New("invalid argument: uri is empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[uri]; ok {
		return r.m[uri], nil
	}
	return Message{}, errors.New(fmt.Sprintf("invalid argument: uri not found [%v]", uri))
}

func (r *messageCache) Uri() []string {
	var uri []string
	r.mu.RLock()
	defer r.mu.RUnlock()
	for key, _ := range r.m {
		uri = append(uri, key)
	}
	sort.Strings(uri)
	return uri
}
