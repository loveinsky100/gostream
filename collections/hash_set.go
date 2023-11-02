/**
 * @author leo
 * @date 2020/8/18 5:51 下午
 */
package collections

import "sync"

const EMPTY = 0

type _HashSet[K any] struct {
	hashMap map[interface{}]int
}

func NewHashSet[K any]() Set[K] {
	return &_HashSet[K]{
		hashMap: map[interface{}]int{},
	}
}

func (hashSet *_HashSet[K]) Add(key K) {
	hashSet.hashMap[key] = EMPTY
}

func (hashSet *_HashSet[K]) Remove(key K) {
	delete(hashSet.hashMap, key)
}

func (hashSet *_HashSet[K]) Contains(key K) bool {
	_, ok := hashSet.hashMap[key]
	return ok
}

func (hashSet *_HashSet[K]) Foreach(consumer func(key K) bool) {
	for k, _ := range hashSet.hashMap {
		_k, ok := k.(K)
		if !ok {
			continue
		}
		if !consumer(_k) {
			break
		}
	}
}

type _SyncHashSet[K any] struct {
	rw  sync.RWMutex
	set Set[K]
}

func NewSyncHashSet[K any]() Set[K] {
	return &_SyncHashSet[K]{
		set: NewHashSet[K](),
	}
}

func (hashSet *_SyncHashSet[K]) Add(key K) {
	hashSet.rw.Lock()
	defer hashSet.rw.Unlock()
	hashSet.set.Add(key)
}

func (hashSet *_SyncHashSet[K]) Remove(key K) {
	hashSet.rw.Lock()
	defer hashSet.rw.Unlock()
	hashSet.set.Remove(key)
}

func (hashSet *_SyncHashSet[K]) Contains(key K) bool {
	hashSet.rw.RLock()
	defer hashSet.rw.RUnlock()
	return hashSet.set.Contains(key)
}

func (hashSet *_SyncHashSet[K]) Foreach(consumer func(key K) bool) {
	hashSet.rw.RLock()
	defer hashSet.rw.RUnlock()
	hashSet.set.Foreach(consumer)
}
