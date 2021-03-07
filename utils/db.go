package utils

import (
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

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

func RedisCache(redisDsn string, redisPassword string, redisDb int) *redis.Client{
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisDsn, // use default Addr
		Password: redisPassword,               // no password set
		DB:       redisDb,                // use default DB
	})
	return redisClient
}