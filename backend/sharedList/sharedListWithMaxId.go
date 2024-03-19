package sharedlist


type MaxId struct {
	Id uint `json:"id"`
}

type EntityWithMaxId interface {
	SetId(MaxId)
}

type sharedListWithMaxId[T EntityWithMaxId] struct {
	sharedList[T]
	maxId MaxId
}

func GetSharedListWithMaxId[T EntityWithMaxId](tablename string, maxId MaxId) *sharedListWithMaxId[T] {
  return &sharedListWithMaxId[T]{
    maxId: maxId,
    sharedList: sharedList[T]{
      tablename: tablename,
    },
  }
}

func (s *sharedListWithMaxId[T]) Add(item T) uint {
	s.lock.Lock()
	defer s.lock.Unlock()

  s.maxId.Id++
  item.SetId(s.maxId)
	s.items = append(s.items, item)
  return s.maxId.Id
}

