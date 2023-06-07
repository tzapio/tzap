package types

type DBCollectionInterface[T any] interface {
	BatchSet(pairs []KeyValue[T]) (int, error)
	Get(key string) (T, bool)
	GetAll() []KeyValue[T]
	ScanGet(key string) (KeyValue[T], bool)
	Set(key string, value T) error
	StartInit()
}
