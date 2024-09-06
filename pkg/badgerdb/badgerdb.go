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

type KVInt struct {
	Key   string
	Value int
}

type KVInterface struct {
	Key   string
	Value int
}

func NewBadgerDB(path string) *BadgerDB {
	opts := badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR)
	opts.ValueThreshold = 16
	opts.ValueLogFileSize = 8 << 20
	opts.NumCompactors = 2
	opts.Compression = options.ZSTD
	opts.MemTableSize = 64 << 20
	opts.BlockCacheSize = 64 << 20
	opts.BloomFalsePositive = 0.01
	opts.NumLevelZeroTables = 2
	opts.NumLevelZeroTablesStall = 10
	opts.NumMemtables = 5
	opts.SyncWrites = false
	opts.ValueLogMaxEntries = 1000000
	opts.MaxLevels = 4

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

func (b *BadgerDB) SetIntegers(ttl time.Duration, kvIntegers ...KVInt) error {
	return b.DB.Update(func(txn *badger.Txn) error {
		for _, kvInt := range kvIntegers {
			valueBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(valueBytes, uint64(kvInt.Value))

			entry := badger.NewEntry([]byte(kvInt.Key), valueBytes).WithTTL(ttl)
			err := txn.SetEntry(entry)
			if err != nil {
				return err
			}
		}
		return nil
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
