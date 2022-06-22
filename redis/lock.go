package rlock

import (
    "errors"
    "github.com/go-redis/redis"
    "github.com/google/uuid"
    "time"
)

var REDIS_KEY = "github.com/czasg/go-lock"

type Unlock func()

func NewRedisLock() *RedisLock {
    return &RedisLock{}
}

type RedisLock struct{}

func (r *RedisLock) Lock() (Unlock, error) {
    uid, err := uuid.NewUUID()
    if err != nil {
        return nil, err
    }
    lock, err := redis.NewClient(nil).SetNX(REDIS_KEY, uid.String(), time.Hour).Result()
    if err != nil {
        return nil, err
    }
    if !lock {
        return nil, errors.New("")
    }
    unlock := func() {
        value, err := redis.NewClient(nil).Get(REDIS_KEY).Result()
        if err != nil {
            return
        }
        if value != uid.String() {
            return
        }
        redis.NewClient(nil).Del(REDIS_KEY)
    }
    return unlock, nil
}
