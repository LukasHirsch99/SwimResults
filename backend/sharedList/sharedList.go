package sharedlist

import (
	"sync"
)

type sharedList[T any] struct {
	lock      sync.RWMutex
	tablename string
	items     []T
}

func GetSharedList[T any](tablename string) *sharedList[T] {
	return &sharedList[T]{
		tablename: tablename,
	}
}

func (s *sharedList[T]) Add(item T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = append(s.items, item)
}

func (s *sharedList[T]) GetTableName() string {
	return s.tablename
}

func (s *sharedList[T]) GetItems() any {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.items
}

func (s *sharedList[T]) GetItemCnt() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

func (s *sharedList[T]) GetItemList() []T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.items
}
