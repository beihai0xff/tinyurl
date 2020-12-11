package storage

import (
	"errors"
	"fmt"
	"github.com/wingsxdu/tinyurl/util"
	"log"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (

	// defaultInitialMmapSize is the initial size of the mmapped region. Setting this larger than
	// the potential max db size can prevent writer from blocking reader.
	// This only works for linux.
	defaultInitialMmapSize = 256 * 1024 * 1024 // 256 MB

	defaultBatchLimit    = 100
	defaultBatchInterval = 500 * time.Millisecond

	StartAt = 123456
)

type Storage interface {
	// Methods to manage key/value pairs
	// Read-Only transactions
	View(bucket, key []byte) ([]byte, error)
	// Create or update a key
	// If the key already exist it will return
	Create(bucket, key, value []byte) error
	// Create or update a key
	// If the key not exist it will Create the key
	Update(bucket, key, value []byte) error
	// Delete a key from bucket
	Delete(bucket, key []byte) error
	// Generate a index for the url and store it
	Index(value []byte) (uint64, error)
	//
	// Batch(func(tx *bolt.Tx) error)

	// Methods to manage a Bucket
	CreateBucket(bucket []byte) error
	DeleteBucket(bucket []byte) error
}

// TODO(beihai): BatchTx
type storage struct {
	db *bolt.DB

	// 两次批量事务提交的最大时间差
	batchInterval time.Duration
	// 指定了一次读写事务中最大的操作数，当超过该阈值时，当前的读写事务会自动提交
	batchLimit int

	buckets map[*bolt.Bucket]time.Duration

	stopc chan struct{}
	donec chan struct{}
}

type Config struct {
	// Path is the file path to the storage file.
	Path string
	// BatchInterval is the maximum time before flushing the BatchTx.
	// 提交两次批量事务的最大时间差，默认 100ms
	BatchInterval time.Duration
	// BatchLimit is the maximum puts before flushing the BatchTx.
	// 指定每个批量读写事务能包含的最多操作个数，当超过这个阈值后，当前批量读写事务会自动提交
	BatchLimit int
	MmapSize   int
	//
	CreateNewFile bool
}

func DefaultConfig() *Config {
	return &Config{
		Path:          "./database/storage.db",
		BatchInterval: defaultBatchInterval,
		BatchLimit:    defaultBatchLimit,
		MmapSize:      defaultInitialMmapSize,
	}
}

func New(c *Config) (Storage, error) {
	return newStorage(c)
}

func newStorage(c *Config) (Storage, error) {
	err := os.Mkdir("./database/", 0600)
	// Open the ./tinyUrl/storage.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(c.Path, 0600, &bolt.Options{Timeout: 3 * time.Second, InitialMmapSize: c.MmapSize})
	if err != nil {
		log.Panicln(err)
	}
	s := &storage{
		db: db,

		batchInterval: c.BatchInterval,
		batchLimit:    c.BatchLimit,

		stopc: make(chan struct{}),
		donec: make(chan struct{}),
	}
	exist, err := s.tryCreateBucket([]byte("index"), true)
	if exist && c.CreateNewFile {
		return nil, errors.New("Bucket Already Exist")
	}
	return s, nil
}

// View a k/v pairs in Read-Only transactions.
func (s *storage) View(bucket, key []byte) ([]byte, error) {
	var v []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}
		v = b.Get(key)
		return nil
	})
	// if the key not exist will return nil
	return v, err
}

// Update a key from given bucket.
func (s *storage) Create(bucket, key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}
		if v := b.Get(key); v != nil {
			return ErrKeyAlreadyExist
		}
		return b.Put(key, value)
	})
}

// Update or Create a key from given bucket.
func (s *storage) Update(bucket, key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}

// Delete a key from given bucket.
func (s *storage) Delete(bucket, key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}
		return b.Delete(key)
	})
}

// Index() will Generate a index for the given url(value) and store it.
func (s *storage) Index(value []byte) (uint64, error) {
	var index uint64
	return index, s.db.Update(func(tx *bolt.Tx) error {
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("index"))
		if b == nil {
			return ErrBucketNotFound
		}

		// Generate index for the url.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		index, _ = b.NextSequence()

		// Persist bytes to url bucket.
		return b.Put(util.Utob(index), value)
	})
}

// return a bucket
func (s *storage) CreateBucket(bucket []byte) error {
	_, err := s.tryCreateBucket(bucket, false)
	return err
}

// tryCreateBucket() will create a Bucket if it not exists
// the field exist tell the caller whether the Bucket already exists.
func (s *storage) tryCreateBucket(bucket []byte, start bool) (bool, error) {
	var exist bool
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			// Bucket not exist, create a new Bucket
			exist = false
			b, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return fmt.Errorf("create bucket %s failed: %s", bucket, err)
			}
			if start {
				err = b.SetSequence(StartAt)
				if err != nil {
					log.Panicln(err)
				}
			}
		} else {
			exist = true // Bucket exists
		}
		return nil
	})
	return exist, err
}

// Delete the given bucket
func (s *storage) DeleteBucket(bucket []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bucket)
	})
	if err != nil {
		if err == ErrBucketNotFound {
			err = nil
		}
	}

	return err
}
