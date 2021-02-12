package lru

import "container/list"

type LRUCache struct {
	Cap int
	// Element 中存储 pair 结构体，用哈希表指向链表的节点
	// LRUCache 执行删除操作的时候，需要维护 2 个数据结构，一个是 map，一个是双向链表。
	// 在双向链表中删除淘汰出去的 value，在 map 中删除淘汰出去 value 对应的 key。
	// 如果在双向链表的 value 中不存储 key，那么再删除 map 中的 key 的时候有点麻烦。
	Keys map[uint64]*list.Element
	List *list.List
}

type pair struct {
	K uint64
	V []byte
}

// Input LRUCache Cap
func New(capacity int) *LRUCache {
	return &LRUCache{
		Cap:  capacity,
		Keys: make(map[uint64]*list.Element),
		List: list.New(),
	}
}

// 获取数据
func (c *LRUCache) Get(key uint64) []byte {
	// 在 map 中直接读取双向链表的结点
	if el, ok := c.Keys[key]; ok {
		// 如果 map 中存在，将它移动到双向链表的表头
		c.List.MoveToFront(el)
		return el.Value.(pair).V
	}
	return nil
}

// 插入或更新数据
func (c *LRUCache) Put(key uint64, value []byte) {
	// 先查询 map 中是否存在 key
	if el, ok := c.Keys[key]; ok {
		// 如果存在，更新它的 value
		el.Value = pair{K: key, V: value}
		// 并且把该结点移到双向链表的表头
		c.List.MoveToFront(el)
	} else {
		// 插入 LRU 中
		el := c.List.PushFront(pair{K: key, V: value})
		c.Keys[key] = el
	}
	// 判断是否需要淘汰最后一个结点
	if c.List.Len() > c.Cap {
		el := c.List.Back()
		c.List.Remove(el)
		delete(c.Keys, el.Value.(pair).K)
	}
}

// 删除数据
func (c *LRUCache) Delete(key uint64) {
	// 先查询 map 中是否存在 key
	if el, ok := c.Keys[key]; ok {
		// 如果存在，则删除节点
		c.List.Remove(el)
		delete(c.Keys, el.Value.(pair).K)
	}
}
