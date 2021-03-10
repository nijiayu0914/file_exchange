package main

import (
	"file_exchange/datamodels"
	"file_exchange/routes"
	"file_exchange/services"
	"file_exchange/utils"
	"fmt"
	"github.com/kataras/iris/v12"
)

func main()  {
	app := iris.New()
	var configuration iris.Configuration
	if utils.FileExist("./config/config.yml"){
		fmt.Println("prod mode")
		configuration = iris.YAML("./config/config.yml")
	}else{
		fmt.Println("dev mode")
		configuration = iris.YAML("./config/dev/config.yml")
	}
	otherConfig := configuration.GetOther()
	dbType := otherConfig["DatabaseType"]
	dbDsn := otherConfig["DatabaseDsn"]
	redisDsn := otherConfig["RedisDsn"]
	redisDb := otherConfig["RedisDb"]
	redisPassword := otherConfig["RedisPassword"]
	oSSEndpoint := otherConfig["OSSEndpoint"]
	oSSAccessKeyID := otherConfig["OSSAccessKeyID"]
	oSSAccessKeySecret := otherConfig["OSSAccessKeySecret"]
	oSSBucketName := otherConfig["OSSBucketName"]
	oSSRegionId := otherConfig["OSSRegionId"]
	oSSRamAccessKeyID := otherConfig["OSSRamAccessKeyID"]
	oSSRamAccessKeySecret := otherConfig["OSSRamAccessKeySecret"]
	oSSRoleArn := otherConfig["OSSRoleArn"]
	oSSRoleSessionName := otherConfig["OSSRoleSessionName"]

	db, err := utils.Db(dbType, dbDsn)
	if err != nil {
		panic("failed to connect database")
	}else{
		db.AutoMigrate(&datamodels.User{}, &datamodels.File{})
	}

	redisClient := utils.RedisCache(
		redisDsn.(string), redisPassword.(string), redisDb.(int))

	ossOperator := services.OssOperator{
		Endpoint:        oSSEndpoint.(string),
		AccessKeyId:     oSSAccessKeyID.(string),
		AccessKeySecret: oSSAccessKeySecret.(string),
		OSSBucketName: oSSBucketName.(string),
		OSSRegionId: oSSRegionId.(string),
		OSSRamAccessKeyID: oSSRamAccessKeyID.(string),
		OSSRamAccessKeySecret: oSSRamAccessKeySecret.(string),
		OSSRoleArn: oSSRoleArn.(string),
		OSSRoleSessionName: oSSRoleSessionName.(string),
	}
	err = ossOperator.GetClient()
	if err != nil{
		panic("failed to connect oss")
	}

	routes.Routes(app, db, redisClient, &ossOperator, otherConfig)

	app.Run(
		iris.Addr("0.0.0.0:8080"),
		iris.WithConfiguration(configuration))
}
