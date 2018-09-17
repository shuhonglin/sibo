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
	MYSQL_TYPE byte = 1
	REDIS_TYPE      = 2
)

var (
	DB_TYPE  byte // 最低位为mysql,次低位为redis(eg: 只开启mysql DB_TYPE:01, 只开启redis DB_TYPE:10, 两者都开启 DB_TYPE:11)。
	SQL_DB   *sqlx.DB
	REDIS_DB *redis.Pool
)

type SqlDbParam struct {
	driverName string
	host       string
	port       int
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
	authPswd    string
}

func InitDB() {
	var err error
	DB_TYPE = 1
	if DB_TYPE == MYSQL_TYPE {
		// 从配置文件中读取
		sqlDbParam := &SqlDbParam{
			driverName: "mysql",
			host:       "127.0.0.1",
			port:       3306,
			urlParam:   "",
			user:       "hackway",
			password:   "0663",
			db:         "test",
		}
		//SQL_DB, err = sqlx.Open("mysql", "linshuhong:feiyin@tcp(127.0.0.1:3306)/fairy_cms?charset=utf8")
		SQL_DB, err = sqlx.Open(sqlDbParam.driverName, sqlDbParam.user+":"+sqlDbParam.password+"@tcp("+sqlDbParam.host+":"+strconv.Itoa(sqlDbParam.port)+")/"+sqlDbParam.db+"?"+sqlDbParam.urlParam)
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
			authPswd:"mac_redis",
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
				if redisDbParam.authPswd != "" {
					if _, err = c.Do("AUTH", redisDbParam.authPswd); err != nil {
						c.Close()
						return nil, err
					}
				}

				// 选择db
				_, err = c.Do("SELECT", redisDbParam.db)
				if err != nil {
					c.Close()
					return nil, err
				}
				return c, nil
			},
		}
		cn := REDIS_DB.Get()
		if cn == nil {
			log.Fatalln(err)
		}
	}

}
