package cache

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type TourCache struct {
	mainCache *safeCache
	getter    Getter
}

func NewTourCache(getter Getter, cache Cache) *TourCache {
	return &TourCache{
		mainCache: newSafeCache(cache),
		getter:    getter,
	}
}

func (t *TourCache) Get(key string) interface{} {
	val := t.mainCache.get(key)
	if val != nil {
		return val
	}

	if t.getter != nil {
		val = t.getter.Get(key)
		if val == nil {
			return nil
		}
		t.mainCache.set(key, val)
		return val
	}

	return nil
}

func (t *TourCache) Set(key string, val interface{}) {
	if val == nil {
		return
	}
	t.mainCache.set(key, val)
}

func (t *TourCache) Stat() *Stat {
	return t.mainCache.stat()
}
