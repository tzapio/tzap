package localdb

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	if strings.HasPrefix(db.filePath, "@MEMORY") {
		return nil
	}
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
			if err.Error() == "EOF" {
				break
			}
			println("Old version of .tzap-data/*.db present. Please delete the .tzap-data/*.db ")
			println(err.Error())
			os.Exit(1)
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
	var writer *gobber.GobWriter
	var file *os.File
	if !strings.HasPrefix(db.filePath, "@MEMORY") {
		filer, err := os.OpenFile(db.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		file = filer
		defer file.Close()

		writer = gobber.NewGobWriterIO(file)
	}

	c := 0
	for _, kv := range pairs {
		if reflect.DeepEqual(db.data[kv.Key], kv.Value) {
			tl.Logger.Println("BatchSet - WARNING: Key already exists. ", kv.Key)
			continue
		}
		c++
		if writer != nil {
			if err := writer.Write(kv); err != nil {
				return c, err
			}
		}

		db.scanKeyList = append(db.scanKeyList, kv)
		if !reflectutil.IsZero(kv.Value) {
			db.data[kv.Key] = kv.Value
		} else {
			delete(db.data, kv.Key)
		}
	}

	if writer != nil {
		if err := file.Sync(); err != nil {
			return c, err
		}
	}

	return c, nil
}
