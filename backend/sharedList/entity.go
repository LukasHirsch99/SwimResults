package sharedlist

type Entity interface {
	GetItems() []any
	GetTableName() string
  GetItemCnt() int
  GetColumnNames() []string
  GetRows() [][]any
}
