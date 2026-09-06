package session

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrSessionNoData      = errors.New("session no data")
	ErrRemoveSessionFail  = errors.New("remove session fail")
	ErrMigrateSessionFail = errors.New("migrate session fail")
)

// Storage global session data store interface.
// You can customize the storage medium by implementing this interface.
type Storage interface {
	// Read data from store
	Read(s *Session) (err error)
	// Write data to storage
	Write(s *Session) (err error)
	// Remove data from storage
	Remove(s *Session) (err error)
}

// timeout manager
type tm map[string]*time.Timer

// MemoryStore Local memory storage.
type MemoryStore struct {
	tm
	rw           sync.RWMutex
	store        map[string]*Session
	garbageTruck chan string
}

// NewRAM return local memory storage.
func NewRAM() *MemoryStore {
	s := &MemoryStore{
		rw:           sync.RWMutex{},
		store:        make(map[string]*Session),
		tm:           make(map[string]*time.Timer, 1024),
		garbageTruck: make(chan string, 1024),
	}
	go s.gc()
	return s
}

func (ram *MemoryStore) Read(s *Session) (err error) {
	ram.rw.RLock()
	defer func() {
		ram.rw.RUnlock()
	}()
	if session, ok := ram.store[s.id]; ok {
		s.Values = session.Values
		s.CreateTime = session.CreateTime
		s.ExpireTime = session.ExpireTime
		return nil
	}
	return ErrSessionNoData
}

func (ram *MemoryStore) Write(s *Session) (err error) {
	ram.rw.Lock()
	defer ram.rw.Unlock()
	ram.store[s.id] = s

	if ram.tm[s.id] == nil {
		go func() {
			ram.tm[s.id] = time.NewTimer(time.Until(s.ExpireTime))
			<-ram.tm[s.id].C
			ram.garbageTruck <- s.id
			ram.tm[s.id].Stop()
		}()
	}

	return nil
}

func (ram *MemoryStore) Remove(s *Session) (err error) {
	ram.rw.Lock()
	defer ram.rw.Unlock()
	delete(ram.tm, s.id)
	delete(ram.store, s.id)
	return nil
}

// gc is ram store garbage collection.
func (ram *MemoryStore) gc() {
	for {
		select {
		case sid := <-ram.garbageTruck:
			ram.rw.Lock()
			delete(ram.store, sid)
			ram.rw.Unlock()
		default:

		}
	}
}

// formatPrefix format redis key prefix
func formatPrefix(sid string) string {
	return fmt.Sprintf("%s:%s", "v", sid)
}

// timeoutCtx redis connect timeout
func timeoutCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
}

// expire redis key expire
func expire(t time.Time) time.Duration {
	return time.Until(t)
}

type FileStore struct{}

func (fs FileStore) Read(s *Session) (err error) {
	panic("implement me")
}

func (fs FileStore) Write(s *Session) (err error) {
	panic("implement me")
}

func (fs FileStore) Remove(s *Session) (err error) {
	panic("implement me")
}
