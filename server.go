package tinyurl

import (
	"errors"
	"github.com/wingsxdu/tinyurl/lru"

	"github.com/wingsxdu/tinyurl/base36"
	"github.com/wingsxdu/tinyurl/storage"
	"github.com/wingsxdu/tinyurl/util"
)

var (
	ErrTinyUrlTooSmall = errors.New("the tinyUrl is too small")
	ErrTinyUrlNotExist = errors.New("the tinyUrl is not exist")
	s                  storage.Storage
)

// Don't use init()

var LRU *lru.LRUCache

// New() will init a Storage interface
// if New() return a error, the program could not run, so we should panic
func New() error {
	c := storage.DefaultConfig()
	var err error
	s, err = storage.New(c)
	LRU = lru.New(1000)
	return err
}

// Get() will get a url by tinyUrl
func Get(tinyUrl string) ([]byte, error) {
	index := base36.Decode(tinyUrl)
	if v := LRU.Get(index); v != nil {
		return v, nil
	} else {
		v, err := s.View([]byte("index"), util.Utob(index))
		if err != nil {
			return nil, err
		}
		LRU.Put(index, v)
		return v, err
	}
}

// Create() will create a tinyUrl
func Create(url []byte) (string, error) {
	index, err := s.Index(url)
	if err != nil {
		return "", err
	}
	return base36.Encode(index), nil
}

// Update() will update a tinyUrl, and the tinyUrl will direct to a new url
func Update(tinyUrl, newUrl string) error {
	index := base36.Decode(tinyUrl)
	if index <= storage.StartAt {
		return ErrTinyUrlTooSmall
	}
	oldUrl, err := s.View([]byte("index"), util.Utob(index))
	if err != nil {
		return err
	} else if oldUrl == nil {
		return ErrTinyUrlNotExist
	}
	LRU.Put(index, []byte(newUrl))
	return s.Update([]byte("index"), util.Utob(index), []byte(newUrl))
}

// Update() will delete a tinyUrl from storage, and it will not be found
func Delete(tinyUrl string) error {
	index := base36.Decode(tinyUrl)
	// if err != nil, the delete function is successful
	LRU.Delete(index)
	return s.Delete([]byte("index"), util.Utob(index))
}
