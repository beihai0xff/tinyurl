package server

import (
	"github.com/wingsxdu/tinyurl/base36"
	"github.com/wingsxdu/tinyurl/storage"
	"github.com/wingsxdu/tinyurl/util"
)

var s storage.Storage

// Don't use init()
func InitServer() {
	c := storage.DefaultConfig()
	s = storage.New(c)
}

func GetTinyUrl(tinyUrl string) ([]byte, error) {
	index := base36.Decode(tinyUrl)
	return s.View([]byte("index"), util.Utob(index))
}

func PostTinyUrl(url []byte) (string, error) {
	index, err := s.Index(url)
	if err != nil {
		return "", err
	}
	return base36.Encode(index), nil
}

func PutTinyUrl(tinyUrl, newUrl string) error {
	index := base36.Decode(tinyUrl)
	return s.Update([]byte("index"), util.Utob(index), []byte(newUrl))
}

func DeleteTinyUrl(tinyUrl string) error {
	index := base36.Decode(tinyUrl)
	// if err != nil, the delete function is successful
	return s.Delete([]byte("index"), util.Utob(index))
}
