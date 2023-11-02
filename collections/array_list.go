/**
 * @author leo
 * @date 2020/8/18 5:11 下午
 */
package collections

import (
	"encoding/json"
	"sync"
)

type _ArrayList[T any] struct {
	items []T
}

func NewArrayList[T any](args ...T) List[T] {
	if len(args) > 0 {
		return &_ArrayList[T]{
			items: args,
		}
	}

	return &_ArrayList[T]{
		items: make([]T, 0),
	}
}

func (arrayList *_ArrayList[T]) Add(value T) {
	arrayList.items = append(arrayList.items, value)
}

func (arrayList *_ArrayList[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(arrayList.items) {
		var zero T
		return zero, false
	}

	return arrayList.items[index], true
}

func (arrayList *_ArrayList[T]) Remove(index int) bool {
	if index < 0 || index >= len(arrayList.items) {
		return false
	}

	arrayList.items = append(arrayList.items[:index], arrayList.items[index+1:]...)
	return true
}

func (arrayList *_ArrayList[T]) Foreach(consumer func(value T, index int) bool) {
	for i, v := range arrayList.items {
		r := consumer(v, i)
		if !r {
			break
		}
	}
}

func (arrayList *_ArrayList[T]) Size() int {
	return len(arrayList.items)
}

func (arrayList *_ArrayList[T]) Data() []T {
	return arrayList.items
}

func (arrayList *_ArrayList[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(arrayList.items)
}

func (arrayList *_ArrayList[T]) UnmarshalJSON(data []byte) error {
	var item []T
	err := json.Unmarshal(data, &item)
	if nil != err {
		return err
	}

	arrayList.items = item
	return nil
}

type _SyncArrayList[T any] struct {
	ArrayList List[T]
	rw        sync.RWMutex
}

func NewSyncArrayList[T any]() List[T] {
	return &_SyncArrayList[T]{
		ArrayList: NewArrayList[T](),
	}
}

func NewSyncList[T any](list List[T]) List[T] {
	return &_SyncArrayList[T]{
		ArrayList: list,
	}
}

func (arrayList *_SyncArrayList[T]) Add(value T) {
	arrayList.rw.Lock()
	defer arrayList.rw.Unlock()

	arrayList.ArrayList.Add(value)
}

func (arrayList *_SyncArrayList[T]) Get(index int) (T, bool) {
	arrayList.rw.RLock()
	defer arrayList.rw.RUnlock()

	return arrayList.ArrayList.Get(index)
}

func (arrayList *_SyncArrayList[T]) Remove(index int) bool {
	arrayList.rw.Lock()
	defer arrayList.rw.Unlock()

	return arrayList.ArrayList.Remove(index)
}

func (arrayList *_SyncArrayList[T]) Foreach(consumer func(value T, index int) bool) {
	arrayList.rw.RLock()
	defer arrayList.rw.RUnlock()

	arrayList.ArrayList.Foreach(consumer)
}

func (arrayList *_SyncArrayList[T]) Size() int {
	arrayList.rw.RLock()
	defer arrayList.rw.RUnlock()
	return arrayList.ArrayList.Size()
}

func (arrayList *_SyncArrayList[T]) Data() []T {
	arrayList.rw.RLock()
	defer arrayList.rw.RUnlock()
	return arrayList.ArrayList.Data()
}

func (arrayList *_SyncArrayList[T]) MarshalJSON() ([]byte, error) {
	arrayList.rw.RLock()
	defer arrayList.rw.RUnlock()
	return json.Marshal(arrayList.ArrayList)
}

func (arrayList *_SyncArrayList[T]) UnmarshalJSON(data []byte) error {
	arrayList.rw.Lock()
	defer arrayList.rw.Unlock()
	return json.Unmarshal(data, &arrayList.ArrayList)
}
