package redisImp

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
)

type RuntimeStorageRedisImp struct {
	//TODO 读写锁?
	sync.RWMutex
	RedisPool *redis.Pool
}

func NewRuntimeStorageRedisImp(redisPool *redis.Pool) *RuntimeStorageRedisImp {
	return &RuntimeStorageRedisImp{RedisPool: redisPool}
}

// Redis Method
func (r *RuntimeStorageRedisImp) getRedisConn() (redis.Conn, error) {
	conn := r.RedisPool.Get()
	if err := conn.Err(); err != nil {
		log.Println("Get Redis Connection failed: ", err)
		return conn, err
	}
	return conn, nil
}

func (r *RuntimeStorageRedisImp) GetAllKeys() ([]string, error) {
	conn, err := r.getRedisConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return redis.Strings(conn.Do("keys", "*"))
}

func (r *RuntimeStorageRedisImp) GetValueByKey(key string) (interface{}, error) {
	conn, err := r.getRedisConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Do("GET", key)
}

func (r *RuntimeStorageRedisImp) GetValueByKeyToBytes(key string) ([]byte, error) {
	return redis.Bytes(r.GetValueByKey(key))
}

func (r *RuntimeStorageRedisImp) GetValueByKeyToString(key string) (string, error) {
	return redis.String(r.GetValueByKey(key))
}

func (r *RuntimeStorageRedisImp) GetValueByKeyToBool(key string) (bool, error) {
	return redis.Bool(r.GetValueByKey(key))
}

func (r *RuntimeStorageRedisImp) GetValueByKeyToInt(key string) (int, error) {
	return redis.Int(r.GetValueByKey(key))
}

func (r *RuntimeStorageRedisImp) SetValueByKey(key string, value interface{}) error {
	conn, err := r.getRedisConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	return err
}

func (r *RuntimeStorageRedisImp) DelValueByKey(key string) error {
	conn, err := r.getRedisConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return err
}

func (r *RuntimeStorageRedisImp) CheckKeyExit(key string) (bool, error) {
	conn, err := r.getRedisConn()
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}
