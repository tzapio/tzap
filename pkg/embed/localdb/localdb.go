package localdb

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
)

type FileDB struct {
	filePath    string
	data        map[string]string
	scanKeyList []types.KeyValue
	lock        sync.RWMutex
}

func NewFileDB(filePath string) (*FileDB, error) {
	db := &FileDB{
		filePath: filePath,
		data:     make(map[string]string),
	}

	if err := db.load(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *FileDB) load() error {
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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var kv types.KeyValue
		if err := json.Unmarshal(scanner.Bytes(), &kv); err != nil {
			return err
		}
		db.scanKeyList = append(db.scanKeyList, kv)
		if kv.Value == "" {
			delete(db.data, kv.Key)
			continue
		}
		db.data[kv.Key] = kv.Value
	}
	return scanner.Err()
}
func (db *FileDB) GetAll() []types.KeyValue {
	db.lock.RLock()
	defer db.lock.RUnlock()
	var values []types.KeyValue
	for key, value := range db.data {
		if value == "" {
			continue
		}
		values = append(values, types.KeyValue{Key: key, Value: value})
	}
	return values
}
func (db *FileDB) ScanGet(key string) (types.KeyValue, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	for _, kv := range db.scanKeyList {
		if kv.Key == key {
			return kv, true
		}
	}
	return types.KeyValue{}, false
}
func (db *FileDB) Get(key string) (string, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	value, exists := db.data[key]
	return value, exists
}

func (db *FileDB) Set(key, value string) error {
	_, err := db.BatchSet([]types.KeyValue{{Key: key, Value: value}})
	return err
}

func (db *FileDB) BatchSet(pairs []types.KeyValue) (int, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	file, err := os.OpenFile(db.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	c := 0
	for _, kv := range pairs {
		if db.data[kv.Key] == kv.Value {
			tl.Logger.Println("BatchSet - WARNING: Key already exists. Continue. ", kv.Key)
			continue
		}
		c++
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
		}
		db.scanKeyList = append(db.scanKeyList, kv)

		if kv.Value != "" {
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
