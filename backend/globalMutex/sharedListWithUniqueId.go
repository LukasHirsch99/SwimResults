package globalMutex

import (
	"slices"
)

type sharedListWithUniqueId[T EntityWithId] struct {
  sharedList[T]
  ids []Id
}

func CreateSharedListWithUniqueId() *sharedListWithUniqueId[EntityWithId] {
  return &sharedListWithUniqueId[EntityWithId]{}
}

func (sl *sharedListWithUniqueId[T]) Add(item T) {
  sl.lock.Lock()
  defer sl.lock.Unlock()
  if slices.Contains(sl.ids, item.GetId()) {
    return
  }

  sl.ids = append(sl.ids, item.GetId())
  sl.items = append(sl.items, item)
}
