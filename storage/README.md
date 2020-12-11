## storage 模块文档

#### 设计思想

storage 模块对嵌入式键值数据库 BlotDB 提供的 API 进行了封装，为上层提供了统一的接口实现，并且可以作为一个独立的模块使用。

BlotDB 使用 B+Tree 组织索引，写入数据时会立即落盘存储，并利用`mmap`技术将数据文件映射到内存中，默认映射到内存的文件大小为 256MB（这一设置只在 Linux 平台上有效），BlotDB 非常适合读多写少的场景。

BlotDB 使用 Bucket 存放键值对，其概念类似于`Namespace`，不同 Bucket 间的键值对无法访问。

#### 使用方法

###### 实例化`Storage`接口

你可以使用默认的配置初始化一个`Storage`接口：

```go
package main

import (
   "github.com/wingsxdu/tinyurl/storage"
)

func main() {
	s = storage.New(storage.DefaultConfig())
}
```

或者自定义配置规则：

```go
package main

import (
   "github.com/wingsxdu/tinyurl/storage"
)

func main() {
	s = storage.New(&storage.Config{
		Path:     "./test/storage.db", // 数据文件存储位置
		MmapSize: 1024 * 1024 * 1024,  // 1GB
	})
}
```

当关闭程序并在下次启动时，storage 模块会自动读取数据文件中的数据。

###### 接口

键值对操作相关的接口：

- `View(bucket, key []byte) ([]byte, error)`方法会开启一个只读事务，并查找指定`key`的值，需要注意的是，如果该`key`不存在会返回nil，但是如果`key`而`value`为空会返回`""`；
- `Update(bucket, key, value []byte) error`方法会开启一个读写事务，并查找指定`key`的值，更新或创建一个`key`；
- `Delete(bucket, key []byte) error`方法会删除指定的`key`，这个方法不会对不存在的`key`特殊处理，而是返回一个 nil error；
- Index(value []byte) (uint64, error)`方法在指定的`index`Bucket 中生成一个自增主键，在存储数据后会返回该主键。

Bucket 操作相关的接口：

- `CreateBucket(bucket []byte) error`方法会尝试创建新的 Bucket，如果该 Bucket 已经存在会忽略创建；
- `DeleteBucket(bucket []byte) error`方法会删除指定的 Bucket，如果该 Bucket 不存在会忽略返回的错误。