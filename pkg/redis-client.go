package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	client = &RedisClient{}
)

type RedisClient struct {
	C *redis.Client
}

var RedisClientConn *redis.Client

type LoggerRedis struct {
	Code         string      `json:"code"`
	Timestamp    time.Time   `json:"timestamp"`
	Id           int         `json:"id"`
	Repositories string      `json:"repositories"`
	Column       int         `json:"column"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

func InitializeRedis() *RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379", //6379
		Password:     "",
		DB:           0,                 // Gunakan database default
		DialTimeout:  5 * time.Second,   // Timeout untuk membuat koneksi baru
		ReadTimeout:  30 * time.Second,  // Timeout untuk pembacaan data
		WriteTimeout: 30 * time.Second,  // Timeout untuk penulisan data
		PoolSize:     50,                // Ukuran pool koneksi (jumlah koneksi maksimum yang aktif)
		MinIdleConns: 10,                // Jumlah koneksi idle minimum
		IdleTimeout:  300 * time.Second, // Timeout koneksi idle
	})

	if err := c.Ping().Err(); err != nil {
		fmt.Println("Unable to connect to redis " + err.Error())
		return nil
	}
	client.C = c
	RedisClientConn = c
	fmt.Println("REDIS CONNECTED")
	return client
}

func (client *RedisClient) GetKey(key string, src interface{}) error {
	val, err := client.C.Get(key).Result()
	if err == redis.Nil || err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}

// SetKey set key
func (client *RedisClient) SetKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = client.C.Set(key, cacheEntry, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *RedisClient) DeleteKey(key string) error {

	iter := client.C.Scan(0, key, 0).Iterator()
	for iter.Next() {
		err := client.C.Del(iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func (client *RedisClient) SettexKey(key string, value interface{}, expiration time.Duration) error {
	cacheSettex, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = client.C.SetXX(key, cacheSettex, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (client *RedisClient) ExpireKey(key string, expiration time.Duration) error {
	err := client.C.Expire(key, expiration).Err()
	if err != nil {
		return err
	}
	return nil

}

func (client *RedisClient) FlushAll() error {
	err := client.C.FlushAll().Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *RedisClient) Close() error {
	err := client.C.Close()
	if err != nil {
		return err
	}
	return nil
}

func (client *RedisClient) GetKeyTTL(key string) (time.Duration, error) {
	ttl, err := client.C.TTL(key).Result()
	if err != nil {
		return 0, err
	}
	if ttl < 0 {
		return 0, nil
	}
	return ttl, nil
}
