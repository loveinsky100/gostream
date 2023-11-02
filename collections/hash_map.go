/**
 * @author leo
 * @date 2020/8/18 5:32 下午
 */
package collections

import (
	"sync"
)

type _HashMap[K comparable, V any] struct {
	itemMap map[K]V
}

func NewHashMap[K comparable, V any]() Map[K, V] {
	return &_HashMap[K, V]{
		itemMap: make(map[K]V),
	}
}

func (hashMap *_HashMap[K, V]) Put(key K, value V) {
	hashMap.itemMap[key] = value
}

func (hashMap *_HashMap[K, V]) Remove(key K) {
	delete(hashMap.itemMap, key)
}

func (hashMap *_HashMap[K, V]) Get(key K) (V, bool) {
	v, ok := hashMap.itemMap[key]
	return v, ok
}

func (hashMap *_HashMap[K, V]) Foreach(consumer func(key K, value V) bool) {
	for k, v := range hashMap.itemMap {
		r := consumer(k, v)
		if !r {
			break
		}
	}
}

type _SyncHashMap[K comparable, V any] struct {
	rw      sync.RWMutex
	itemMap Map[K, V]
}

func NewSyncHashMap[K comparable, V any]() Map[K, V] {
	return &_SyncHashMap[K, V]{
		itemMap: NewHashMap[K, V](),
	}
}

func (hashMap *_SyncHashMap[K, V]) Put(key K, value V) {
	hashMap.rw.Lock()
	defer hashMap.rw.Unlock()
	hashMap.itemMap.Put(key, value)
}

func (hashMap *_SyncHashMap[K, V]) Remove(key K) {
	hashMap.rw.Lock()
	defer hashMap.rw.Unlock()
	hashMap.itemMap.Remove(key)
}

func (hashMap *_SyncHashMap[K, V]) Get(key K) (V, bool) {
	hashMap.rw.RLock()
	defer hashMap.rw.RUnlock()
	return hashMap.itemMap.Get(key)
}

func (hashMap *_SyncHashMap[K, V]) Foreach(consumer func(key K, value V) bool) {
	hashMap.rw.RLock()
	defer hashMap.rw.RUnlock()
	hashMap.itemMap.Foreach(consumer)
}
