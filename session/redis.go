package session

import "github.com/gomodule/redigo/redis"

//Redis Method
func (sess *SessionsManager) GetRedisConn() (redis.Conn, error) {
	conn := sess.RedisPool.Get()
	if err := conn.Err(); err != nil {
		return conn, err
	}
	return conn, nil
}

func (sess *SessionsManager) CheckRedisKeyExit(key string) (bool, error) {
	conn, err := sess.GetRedisConn()
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

func (sess *SessionsManager) DelRedisByKey(key string) error {
	conn, err := sess.GetRedisConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return err
}

func (sess *SessionsManager) SetRedisKeyValue(key string, value interface{}) error {
	conn, err := sess.GetRedisConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	return err
}

func (sess *SessionsManager) GetAllRedisKey() ([]string, error) {
	conn, err := sess.GetRedisConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return redis.Strings(conn.Do("keys", "*"))
}

func (sess *SessionsManager) GetRedisValueByKey(key string) (interface{}, error) {
	conn, err := sess.GetRedisConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Do("GET", key)
}

func (sess *SessionsManager) GetRedisValueByKeyToBytes(key string) ([]byte, error) {
	return redis.Bytes(sess.GetRedisValueByKey(key))
}

func (sess *SessionsManager) GetRedisValueByKeyToString(key string) (string, error) {
	return redis.String(sess.GetRedisValueByKey(key))
}

func (sess *SessionsManager) GetRedisValueByKeyToBool(key string) (bool, error) {
	return redis.Bool(sess.GetRedisValueByKey(key))
}

func (sess *SessionsManager) GetRedisValueByKeyToInt(key string) (int, error) {
	return redis.Int(sess.GetRedisValueByKey(key))
}

//Redis Method End
