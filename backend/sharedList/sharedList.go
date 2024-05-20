package sharedlist

import (
	"reflect"
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

func (s *sharedList[T]) GetItems() []any {
	s.lock.RLock()
	defer s.lock.RUnlock()

	r := make([]any, s.GetItemCnt())
	for i, e := range s.items {
		r[i] = e
	}
	return r
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

func (s *sharedList[T]) GetColumnNames() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var zero T
	var entityType = reflect.TypeOf(zero)
	columnNames := make([]string, entityType.NumField())
	for i := 0; i < entityType.NumField(); i++ {
		columnNames[i] = entityType.Field(i).Name
	}
	return columnNames
}

func (s *sharedList[T]) GetRows() [][]any {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var zero T
	var entityType = reflect.TypeOf(zero)
	rows := make([][]any, s.GetItemCnt())

	for i, row := range s.items {
		r := make([]any, entityType.NumField())
		for j := 0; j < entityType.NumField(); j++ {
			if reflect.ValueOf(row).Field(j).Interface() != nil {
				r[j] = reflect.ValueOf(row).Field(j).Interface()
			}
		}
		rows[i] = r
	}
	return rows
}
