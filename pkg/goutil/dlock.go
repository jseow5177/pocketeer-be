package goutil

import "sync"

type DLock struct {
	locks map[string]chan struct{}
	mu    sync.RWMutex
}

func NewDLock() *DLock {
	return &DLock{
		locks: make(map[string]chan struct{}),
	}
}

func (dl *DLock) TryLock(key string) bool {
	dl.mu.Lock()
	ch, ok := dl.locks[key]
	if ok {
		dl.mu.Unlock()
		select {
		case ch <- struct{}{}:
			return true
		default:
			return false
		}
	}

	ch = make(chan struct{}, 1)
	dl.locks[key] = ch
	dl.mu.Unlock()

	ch <- struct{}{}
	return true
}

func (dl *DLock) Unlock(key string) {
	dl.mu.Lock()
	if ch, ok := dl.locks[key]; ok {
		select {
		case <-ch:
			close(ch)
			delete(dl.locks, key)
		default:
		}
	}
	dl.mu.Unlock()
}
