package globalMutex

type sharedListWithMaxId[T any] struct {
  sharedList[EntityWithAutoId]
  maxId AutoId
}

func CreateSharedListWithMaxId() *sharedListWithMaxId[EntityWithAutoId] {
  return &sharedListWithMaxId[EntityWithAutoId]{}
}

func (sl *sharedListWithMaxId[T]) Add(item EntityWithAutoId) {
  sl.lock.Lock()
  defer sl.lock.Unlock()
  sl.maxId.Id++
  item.SetId(sl.maxId)
  sl.items = append(sl.items, item)
}

