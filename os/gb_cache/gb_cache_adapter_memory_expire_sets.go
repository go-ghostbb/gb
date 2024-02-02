package gbcache

import (
	gbset "ghostbb.io/gb/container/gb_set"
	"sync"
)

type adapterMemoryExpireSets struct {
	mu         sync.RWMutex         // expireSetMu ensures the concurrent safety of expireSets map.
	expireSets map[int64]*gbset.Set // expireSets is the expiring timestamp to its key set mapping, which is used for quick indexing and deleting.
}

func newAdapterMemoryExpireSets() *adapterMemoryExpireSets {
	return &adapterMemoryExpireSets{
		expireSets: make(map[int64]*gbset.Set),
	}
}

func (d *adapterMemoryExpireSets) Get(key int64) (result *gbset.Set) {
	d.mu.RLock()
	result = d.expireSets[key]
	d.mu.RUnlock()
	return
}

func (d *adapterMemoryExpireSets) GetOrNew(key int64) (result *gbset.Set) {
	if result = d.Get(key); result != nil {
		return
	}
	d.mu.Lock()
	if es, ok := d.expireSets[key]; ok {
		result = es
	} else {
		result = gbset.New(true)
		d.expireSets[key] = result
	}
	d.mu.Unlock()
	return
}

func (d *adapterMemoryExpireSets) Delete(key int64) {
	d.mu.Lock()
	delete(d.expireSets, key)
	d.mu.Unlock()
}
