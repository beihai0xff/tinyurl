package storage

import (
	"github.com/wingsxdu/tinyurl/util"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

var s Storage

var s1 *storage

func init() {
	os.RemoveAll("./tinyUrl/")
	s = New(&Config{
		Path:          "./tinyUrl/test.db",
		BatchInterval: defaultBatchInterval,
		BatchLimit:    defaultBatchLimit,
		MmapSize:      defaultInitialMmapSize,
	})

	db, err := bolt.Open("./tinyUrl/test2.db", 0600, &bolt.Options{Timeout: 3 * time.Second, InitialMmapSize: defaultInitialMmapSize})
	if err != nil {
		log.Panicln(err)
	}
	s1 = &storage{
		db: db,
	}
}

func Test_storage_Index(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
		want  uint64
	}{
		{"test_short_url", []byte("https://www.google.com/"), 123457},
		{"test_long_url", []byte("https://www.google.com/search?sxsrf=ALeKk00rEgE8Gd7-KSZTZUxVkWSzq6exKw%3A1592901527548&ei=l7_xXrmMIYfd9QOOv6CACg&q=%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE&oq=%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE&gs_lcp=CgZwc3ktYWIQAzIECCMQJzIECAAQHlDIKVj2K2CdLmgAcAB4AIABiQOIAcMHkgEFMi0yLjGYAQCgAQGqAQdnd3Mtd2l6&sclient=psy-ab&ved=0ahUKEwj5s9jNxJfqAhWHbn0KHY4fCKAQ4dUDCAw&uact=5"),
			123458},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Index(tt.value)
			if err != nil {
				t.Errorf("Index() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Index() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_Update(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
		key    []byte
		value  []byte
	}{
		{"test_short_url", []byte("index"), util.Utob(uint64(123457)), []byte("https://cn.bing.com/")},
		{"test_long_url", []byte("index"), util.Utob(uint64(123458)), []byte("https://cn.bing.com/search?q=%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE&qs=n&form=QBLHCN&sp=-1&pq=zhao%27hua%27xi%27shi&sc=3-15&sk=&cvid=4FA5BBC53EA84E6B93A6DEC3F006AA4D")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Update(tt.bucket, tt.key, tt.value); err != nil {
				t.Errorf("Update() error = %v", err)
			}
		})
	}
}

func Test_storage_View(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
		key    []byte
		want   []byte
	}{
		{"test_short_url", []byte("index"), util.Utob(uint64(123457)), []byte("https://cn.bing.com/")},
		{"test_long_url", []byte("index"), util.Utob(uint64(123458)), []byte("https://cn.bing.com/search?q=%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE&qs=n&form=QBLHCN&sp=-1&pq=zhao%27hua%27xi%27shi&sc=3-15&sk=&cvid=4FA5BBC53EA84E6B93A6DEC3F006AA4D")},
	}
	Test_storage_Update(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.View(tt.bucket, tt.key)
			if err != nil {
				t.Errorf("View() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("View() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_Delete(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
		key    []byte
	}{
		{"test_short_url", []byte("index"), util.Utob(uint64(123457))},
		{"test_long_url", []byte("index"), util.Utob(uint64(123458))},
	}
	Test_storage_Update(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Delete(tt.bucket, tt.key); err != nil {
				t.Errorf("Delete() error = %v", err)
			}
		})
	}
}

func Test_storage_tryCreateBucket(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
		exits  bool
	}{
		{"test_CreateBucket", []byte("testCreateBucket"), false},
		{"test_CreateBucket", []byte("testCreateBucket"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s1.tryCreateBucket(tt.bucket, false)
			if err != nil {
				t.Errorf("tryCreateBucket() error = %v", err)
				return
			}
			if got != tt.exits {
				t.Errorf("tryCreateBucket() got = %v, exits %v", got, tt.exits)
			}
		})
	}
}

func Test_storage_CreateBucket(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
	}{
		{"test_CreateBucket", []byte("testCreateBucket")},
		{"test_CreateBucket", []byte("testCreateBucket")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateBucket(tt.bucket)
			if err != nil {
				t.Errorf("CreateBucket() error = %v, ", err)
				return
			}
		})
	}
}

func Test_storage_DeleteBucket(t *testing.T) {
	tests := []struct {
		name    string
		bucket  []byte
		wantErr error
	}{
		{"test_DeleteExistBucket", []byte("ExistBucket"), nil},
		{"test_DeleteNotExistBucket", []byte("notExistBucket"), ErrBucketNotFound},
	}
	err := s.CreateBucket(tests[0].bucket)
	if err != nil {
		t.Errorf("CreateBucket() error = %v, ", err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.DeleteBucket(tt.bucket); err != nil {
				if err == tt.wantErr {
					return
				} else {
					t.Errorf("DeleteBucket() error = %v", err)
				}
			}
		})
	}
}
