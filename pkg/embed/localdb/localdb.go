package localdb

import (
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed/localdb/gobber"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/util/reflectutil"
)

type FileDB[T any] struct {
	filePath    string
	data        map[string]T
	scanKeyList []types.KeyValue[T]
	lock        sync.RWMutex
	ready       chan struct{} // signal when the data is ready
}

func NewFileDB[T any](filePath string) (*FileDB[T], error) {
	tl.Logger.Printf("NewFileDB: %s\n", filePath)
	db := &FileDB[T]{
		filePath: filePath,
		data:     make(map[string]T),
		ready:    make(chan struct{}),
	}
	go func() {
		// close the channel to signal we're done
		defer close(db.ready)

		if err := db.load(); err != nil {
			tl.Logger.Printf("error loading data: %v", err)
			// handle error, you may want to pass it out to the caller
			// or store it in the FileDB structure and check it before usage
		}
	}()

	return db, nil
}
func (db *FileDB[T]) waitReady() {
	<-db.ready
}
func (db *FileDB[T]) load() error {
	file, err := os.Open(db.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			absPath, err := filepath.Abs(db.filePath)
			if err != nil {
				return err
			}
			folderPath := filepath.Dir(absPath)
			if err := os.MkdirAll(folderPath, 0755); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	defer file.Close()
	reader := gobber.NewGobReaderIO(file)
	for {
		var kv types.KeyValue[T]
		err := reader.Read(&kv)
		if err != nil {
			break
		}
		db.scanKeyList = append(db.scanKeyList, kv)
		if reflectutil.IsZero(kv.Value) {
			delete(db.data, kv.Key)
			continue
		}
		db.data[kv.Key] = kv.Value

	}
	tl.Logger.Printf("Done NewFileDB: %s Loaded: %d\n", db.filePath, len(db.data))
	return nil
}
func (db *FileDB[T]) GetAll() []types.KeyValue[T] {
	db.waitReady()
	db.lock.RLock()
	defer db.lock.RUnlock()
	var values []types.KeyValue[T]
	for key, value := range db.data {
		if reflectutil.IsZero(value) {
			continue
		}
		values = append(values, types.KeyValue[T]{Key: key, Value: value})
	}
	return values
}
func (db *FileDB[T]) ScanGet(key string) (types.KeyValue[T], bool) {
	db.waitReady()
	db.lock.RLock()
	defer db.lock.RUnlock()
	for _, kv := range db.scanKeyList {
		if kv.Key == key {
			return kv, true
		}
	}
	return types.KeyValue[T]{}, false
}
func (db *FileDB[T]) Get(key string) (T, bool) {
	db.waitReady()
	db.lock.RLock()
	defer db.lock.RUnlock()
	value, exists := db.data[key]
	return value, exists
}

func (db *FileDB[T]) Set(key string, value T) error {
	_, err := db.BatchSet([]types.KeyValue[T]{{Key: key, Value: value}})
	return err
}
func (db *FileDB[T]) BatchSet(pairs []types.KeyValue[T]) (int, error) {
	db.waitReady()
	db.lock.Lock()
	defer db.lock.Unlock()

	file, err := os.OpenFile(db.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	c := 0
	writer := gobber.NewGobWriterIO(file)
	for _, kv := range pairs {
		if reflect.DeepEqual(db.data[kv.Key], kv.Value) {
			tl.Logger.Println("BatchSet - WARNING: Key already exists. ", kv.Key)
			continue
		}
		c++
		err := writer.Write(kv)
		if err != nil {
			return c, err
		}
		/*
			data, err := json.Marshal(kv)
			if err != nil {
				println("BatchSet - WARNING: Cannot marshal data. Continue. ", err.Error())
				continue
			}

			if _, err := file.Write(data); err != nil {
				return c, err
			}
			if _, err := file.WriteString("\n"); err != nil {
				return c, err
			}*/
		db.scanKeyList = append(db.scanKeyList, kv)

		if !reflectutil.IsZero(kv.Value) {
			db.data[kv.Key] = kv.Value
		} else {
			delete(db.data, kv.Key)
		}
	}
	if err := file.Sync(); err != nil {
		return c, err
	}

	return c, nil
}
