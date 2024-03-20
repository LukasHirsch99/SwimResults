package sharedlist

import (
	"slices"

	"github.com/gocolly/colly"
)

type EntityWithUniqueId interface {
	GetId() uint
}

type sharedListWithUniqueId[T EntityWithUniqueId] struct {
	sharedList[T]
	ids []uint
}

func GetSharedListWithUniqueId[T EntityWithUniqueId](tablename string, ids []uint) *sharedListWithUniqueId[T] {
	return &sharedListWithUniqueId[T]{
		ids: ids,
		sharedList: sharedList[T]{
			tablename: tablename,
		},
	}
}

func GetSharedListWithUniqueIdWithItems[T EntityWithUniqueId](tablename string, items []T) *sharedListWithUniqueId[T] {
	sl := &sharedListWithUniqueId[T]{
		sharedList: sharedList[T]{
			tablename: tablename,
			items:     items,
		},
	}
	for _, item := range items {
		sl.ids = append(sl.ids, item.GetId())
	}
	return sl
}

func (s *sharedListWithUniqueId[T]) addWithoutLock(item T) {
	if slices.Contains(s.ids, item.GetId()) {
		return
	}

	s.items = append(s.items, item)
	s.ids = append(s.ids, item.GetId())
}

func (s *sharedListWithUniqueId[T]) Add(item T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.addWithoutLock(item)
}

func (s *sharedListWithUniqueId[T]) EnsureExists(id uint, row *colly.HTMLElement, creator func(uint, *colly.HTMLElement) T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if slices.Contains(s.ids, id) {
		return
	}

	item := creator(id, row)
	s.addWithoutLock(item)
}

func (s *sharedListWithUniqueId[T]) GetItemById(id uint) T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, item := range s.items {
		if item.GetId() == id {
			return item
		}
	}
	return *new(T)
}

func (s *sharedListWithUniqueId[T]) SetItem(item T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for idx, i := range s.items {
		if i.GetId() == item.GetId() {
			s.items[idx] = item
			return
		}
	}
	s.items = append(s.items, item)
}
