package fast

import (
	"container/list"
	"sync"
)

type cacheShard struct {
	locker sync.RWMutex

	// 最大存放 entry 个数
	maxEntries int
	// 当一个 entry 从缓存中移除是调用该回调函数，默认为 nil
	// groupcache 中的 key 是任意的可比较类型；value 是 interface{}
	onEvicted func(key string, value interface{})

	ll    *list.List
	cache map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

// new 创建一个新的 cacheShard，如果 maxBytes 是 0，表示没有容量限制
func newCacheShard(maxEntries int, onEvicted func(key string, value interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// set 往 Cache 尾部增加一个元素（如果已经存在，则放入尾部，并更新值）
func (c *cacheShard) set(key string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		return
	}

	en := &entry{key, value}
	e := c.ll.PushBack(en)
	c.cache[key] = e

	if c.maxEntries > 0 && c.ll.Len() > c.maxEntries {
		c.removeElement(c.ll.Front())
	}
}

// get 从 cache 中获取 key 对应的值，nil 表示 key 不存在
func (c *cacheShard) get(key string) interface{} {
	c.locker.RLock()
	defer c.locker.RUnlock()

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}

	return nil
}

// del 从 cache 中删除 key 对应的元素
func (c *cacheShard) del(key string) {
	c.locker.Lock()
	defer c.locker.Unlock()

	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

// delOldest 从 cache 中删除最旧的记录
func (c *cacheShard) delOldest() {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.removeElement(c.ll.Front())
}

// len 返回当前 cache 中的记录数
func (c *cacheShard) len() int {
	c.locker.RLock()
	defer c.locker.RUnlock()

	return c.ll.Len()
}

func (c *cacheShard) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	c.ll.Remove(e)
	en := e.Value.(*entry)
	delete(c.cache, en.key)

	if c.onEvicted != nil {
		c.onEvicted(en.key, en.value)
	}
}
