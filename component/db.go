package component

import (
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	MYSQL_TYPE byte = 0
	REDIS_TYPE      = 1
)

var (
	DB_TYPE  byte
	SQL_DB   *sqlx.DB
	REDIS_DB *redis.Pool
)

type SqlDbParam struct {
	driverName string
	host       string
	port       string
	urlParam   string
	user       string
	password   string
	db         string
}

type RedisDbParam struct {
	host        string
	port        int
	maxIdle     int
	maxActive   int
	idleTimeout time.Duration
	db          int
}

func init() {
	var err error
	DB_TYPE = 1
	if DB_TYPE == MYSQL_TYPE {
		// 从配置文件中读取
		sqlDbParam := &SqlDbParam{
			driverName: "mysql",
			host:       "127.0.0.1",
			port:       "3306",
			urlParam:   "",
			user:       "linshuhong",
			password:   "feiyin",
			db:         "fairy_cms",
		}
		//SQL_DB, err = sqlx.Open("mysql", "linshuhong:feiyin@tcp(127.0.0.1:3306)/fairy_cms?charset=utf8")
		SQL_DB, err = sqlx.Open(sqlDbParam.driverName, sqlDbParam.user+":"+sqlDbParam.password+"@tcp("+sqlDbParam.host+":"+sqlDbParam.port+")/"+sqlDbParam.db+"?"+sqlDbParam.urlParam)
		if err != nil {
			log.Fatalln(err)
		}
		err = SQL_DB.Ping()
		if err != nil {
			log.Fatalln(err)
		}
	} else if DB_TYPE == REDIS_TYPE {
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		redisDbParam := &RedisDbParam{
			host:        "127.0.0.1",
			port:        6379,
			maxIdle:     1,
			maxActive:   10,
			idleTimeout: 180 * time.Second,
			db:          1,
		}
		REDIS_DB = &redis.Pool{
			MaxIdle:     redisDbParam.maxIdle,
			MaxActive:   redisDbParam.maxActive,
			IdleTimeout: redisDbParam.idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", redisDbParam.host+":"+strconv.Itoa(redisDbParam.port))
				if err != nil {
					return nil, err
				}
				// 选择db
				c.Do("SELECT", redisDbParam.db)
				return c, nil
			},
		}
	}

}
