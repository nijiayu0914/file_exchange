package utils

import (
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Db 初始化gorm Db对象，支持mysql和postgresql
// param dbType: 数据库类型 mysql,postgresql...
// param dbDsn: 数据库连接
func Db(dbType interface{}, dbDsn interface{})(db *gorm.DB, err error){
	if dbType.(string) == "mysql"{
		db, err = gorm.Open(mysql.Open(dbDsn.(string)), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	}else if dbType.(string) == "postgres"{
		db, err = gorm.Open(postgres.Open(dbDsn.(string)), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	}
	return db, err
}

// Paginate gorm 分页器
// param page: 页码
// param pageSize: 每页大小
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// RedisCache 初始化Redis Client
// param redisDsn 地址
// param redisPassword 数据库密码
// param redisDb 数据库 0, 1, 2......
func RedisCache(redisDsn string, redisPassword string,
	redisDb int) *redis.Client{
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisDsn,
		Password: redisPassword,
		DB:       redisDb,
	})
	return redisClient
}