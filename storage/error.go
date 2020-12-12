package storage

import (
	"errors"
	bolt "go.etcd.io/bbolt"
)

var (
	ErrBucketNotFound     = bolt.ErrBucketNotFound
	ErrKeyAlreadyExist    = errors.New("the key already exist")
	ErrBucketAlreadyExist = errors.New("the bucket already exist")
)
