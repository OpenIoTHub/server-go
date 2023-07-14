// Package memImp 内存实现运行时存储接口
package memImp

import (
	"errors"
	"sync"
)

type RuntimeStorageMemImp struct {
	safeMap sync.Map
}

func NewRuntimeStorageMemImp() *RuntimeStorageMemImp {
	return &RuntimeStorageMemImp{}
}

func (r *RuntimeStorageMemImp) GetAllKeys() ([]string, error) {
	keys := make([]string, 0)
	r.safeMap.Range(func(key, value any) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys, nil
}

func (r *RuntimeStorageMemImp) GetValueByKey(key string) (any, error) {
	if value, ok := r.safeMap.Load(key); ok {
		return value, nil
	}
	return nil, errors.New("key not exist")
}

func (r *RuntimeStorageMemImp) GetValueByKeyToBytes(key string) ([]byte, error) {
	v, err := r.GetValueByKey(key)
	return v.([]byte), err
}

func (r *RuntimeStorageMemImp) GetValueByKeyToString(key string) (string, error) {
	v, err := r.GetValueByKey(key)
	return v.(string), err
}

func (r *RuntimeStorageMemImp) GetValueByKeyToBool(key string) (bool, error) {
	v, err := r.GetValueByKey(key)
	return v.(bool), err
}

func (r *RuntimeStorageMemImp) GetValueByKeyToInt(key string) (int, error) {
	v, err := r.GetValueByKey(key)
	return v.(int), err
}

func (r *RuntimeStorageMemImp) SetValueByKey(key string, value any) error {
	r.safeMap.Store(key, value)
	return nil
}

func (r *RuntimeStorageMemImp) DelValueByKey(key string) error {
	r.safeMap.Delete(key)
	return nil
}

func (r *RuntimeStorageMemImp) CheckKeyExit(key string) (bool, error) {
	if _, ok := r.safeMap.Load(key); ok {
		return ok, nil
	}
	return false, nil
}
