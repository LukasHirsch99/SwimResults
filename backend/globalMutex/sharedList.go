package globalMutex

import (
	"encoding/json"
	"sync"
)

type sharedList[T any] struct {
	lock  sync.Mutex
	items []T
}

func CreateSharedList[T any]() *sharedList[T] {
  return &sharedList[T]{}
}

func (sl *sharedList[T]) Add(item T) {
  sl.lock.Lock()
  defer sl.lock.Unlock()
  sl.items = append(sl.items, item)
}

func (sl *sharedList[T]) MarshalJSON() ([]byte, error) {
  return json.Marshal(sl.items)
}
