package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/gommon/log"
)

var (
	redisConnPoll redis.Pool
)

func InitRedisPool() {

	redisConnPoll = redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Errorf("Error while dialing for redis %v", err)
				return nil, err
			}
			if _, err := conn.Do("SELECT", 0); err != nil {
				conn.Close()
				return nil, err
			}
			//log.Infof("Redis Connection Pool Started!")
			return conn, nil
		},
	}
}

func GetRedisConnFromPool() redis.Conn {
	redisConn := redisConnPoll.Get()
	return redisConn
}
