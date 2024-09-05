package badgerdb

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
	"log"
	"time"
)

type BadgerDB struct {
	DB *badger.DB
}

func NewBadgerDB(path string) *BadgerDB {
	opts := badger.DefaultOptions(path).WithLoggingLevel(badger.INFO)
	opts.ValueThreshold = 16
	opts.ValueLogFileSize = 8 << 20
	opts.NumCompactors = 4
	opts.Compression = options.ZSTD

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	return &BadgerDB{
		DB: db,
	}
}

func (b *BadgerDB) SetInt(key string, value int, ttl time.Duration) error {
	return b.DB.Update(func(txn *badger.Txn) error {
		valueBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(valueBytes, uint64(value))

		entry := badger.NewEntry([]byte(key), valueBytes).WithTTL(ttl)
		return txn.SetEntry(entry)
	})
}

func (b *BadgerDB) GetInt(key string) (int, error) {
	var intValue int
	err := b.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}
		if item != nil {
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			intValue = int(binary.BigEndian.Uint64(val))
		}
		return nil
	})
	return intValue, err
}

func (b *BadgerDB) SetObject(key string, value interface{}, ttl time.Duration) error {
	return b.DB.Update(func(txn *badger.Txn) error {
		valueBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}

		entry := badger.NewEntry([]byte(key), valueBytes).WithTTL(ttl)
		return txn.SetEntry(entry)
	})
}

func (b *BadgerDB) GetObject(key string, target interface{}) error {
	return b.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}
		if item != nil {
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			return json.Unmarshal(val, target)
		}
		return nil
	})
}

func (b *BadgerDB) DeleteKey(key string) error {
	return b.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *BadgerDB) Close() {
	err := b.DB.Close()
	if err != nil {
		log.Fatal(err)
	}
}
