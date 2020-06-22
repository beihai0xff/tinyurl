package storage

import (
	"fmt"
	"log"
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
)

type Storage interface {
	// Methods to manage key/value pairs
	// Read-Only transactions
	View(bucket, key []byte) ([]byte, error)
	// Create or update a key
	Update(bucket, key, value []byte) error
	// delete a key from bucket
	Delete(bucket, key []byte) error
	//
	// Batch(func(tx *bolt.Tx) error)

	// Methods to manage a Bucket
	BucketCreate(bucket []byte) error
	BucketDelete(bucket []byte) error
}

type storage struct {
	db *bolt.DB

	// 两次批量事务提交的最大时间差
	batchInterval time.Duration
	// 指定了一次读写事务中最大的操作数，当超过该阈值时，当前的读写事务会自动提交
	batchLimit int

	stopc chan struct{}
	donec chan struct{}
}

type Config struct {
	// Path is the file path to the backend file.
	Path string
	// BatchInterval is the maximum time before flushing the BatchTx.
	// 提交两次批量事务的最大时间差，默认 100ms
	BatchInterval time.Duration
	// BatchLimit is the maximum puts before flushing the BatchTx.
	// 指定每个批量读写事务能包含的最多操作个数，当超过这个阈值后，当前批量读写事务会自动提交
	BatchLimit int
	MmapSize   int64
}

func DefaultConfig() *Config {
	return &Config{
		Path:          "./tinyUrl/storage.db",
		BatchInterval: defaultBatchInterval,
		BatchLimit:    defaultBatchLimit,
		MmapSize:      defaultInitialMmapSize,
	}
}

func New(c *Config) Storage {
	return newStorage(c)
}

func newStorage(c *Config) Storage {
	// Open the ./tinyUrl/storage.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(c.Path, 0600, &bolt.Options{Timeout: 3 * time.Second, InitialMmapSize: c.MmapSize})
	if err != nil {
		log.Fatal(err)
	}
	s := &storage{
		db: db,

		batchInterval: c.BatchInterval,
		batchLimit:    c.BatchLimit,

		stopc: make(chan struct{}),
		donec: make(chan struct{}),
	}
	// 启动一个单独的协程，其中会定时提交当前的读写事务，并开启新的读写事务
	// go s.run()
	return s
}

func (s *storage) View(bucket, key []byte) ([]byte, error) {
	var v []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return errorBucketNotFound
		}
		v = b.Get(key)
		return nil
	})
	// 不存在的键值对返回 nil
	return v, err
}

func (s *storage) Update(bucket, key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}

func (s *storage) Delete(bucket, key []byte) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return errorBucketNotFound
		}
		return b.Delete(key)
	})
}

func (s *storage) BucketCreate(bucket []byte) (*bolt.Bucket, error) {
	return s.createBucket(bucket)
}

func (s *storage) createBucket(bucket []byte) (*bolt.Bucket, error) {
	var (
		b   *bolt.Bucket
		err error
	)
	err = s.db.Update(func(tx *bolt.Tx) error {
		b, err = tx.CreateBucket(bucket)
		if err != nil {
			return fmt.Errorf("create bucket %s failed: %s", bucket, err)
		}
		return nil
	})
	return b, err
}

func (s *storage) BucketDelete(bucket []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bucket)
	})
}

// func (s *storage) run() {
// 	defer close(s.donec)
// 	// 定时器
// 	t := time.NewTimer(s.batchInterval)
// 	defer t.Stop()
// 	for {
// 		select {
// 		case <-t.C:
// 		case <-s.stopc:
// 		}
// 		t.Reset(s.batchInterval)
// 	}
// }
