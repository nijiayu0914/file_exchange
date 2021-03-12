// main 项目入口文件，配置加载，数据库初始化等
package main

import (
	"file_exchange/datamodels"
	"file_exchange/routes"
	"file_exchange/services"
	"file_exchange/utils"
	"flag"
	"github.com/kataras/iris/v12"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)


// 项目入口
func main()  {
	// 配置文件存储类型，本地文件存储 or URL GET请求，默认文件存储
	configType := flag.String(
		"config_type",
		"file",
		"配置文件获取类型，file为本地文件，http问网络路径")

	// 配置文件存储地址，默认./config/dev/config.yml
	configPath := flag.String(
		"config_path",
		"./config/dev/config.yml",
		"配置文件地址")
	flag.Parse()

	app := iris.New()
	var configuration iris.Configuration
	if *configType == "file"{
		log.Println(*configType)
		log.Println(*configPath)
		if utils.FileExist(*configPath){
			log.Println("config init")
			configuration = iris.YAML(*configPath)
		} else{
			log.Println("config not exist")
		}
	}else if *configType == "http"{
		log.Println(*configType)
		log.Println(*configPath)
		configuration = iris.DefaultConfiguration()
		resp, err := http.Get(*configPath)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		ymlConfig, err := ioutil.ReadAll(resp.Body)
		err = yaml.Unmarshal(ymlConfig, &configuration)
		if err != nil{
			log.Println("config 解析失败")
		}else{
			log.Println("config init")
		}
	}else{
		log.Println("config type is wrong")
	}

	otherConfig := configuration.GetOther()
	port := otherConfig["Port"]
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

	// 初始化数据库
	db, err := utils.Db(dbType, dbDsn)
	if err != nil {
		panic("failed to connect database")
	}else{
		db.AutoMigrate(&datamodels.User{}, &datamodels.File{})
	}

	// 初始化redis
	redisClient := utils.RedisCache(
		redisDsn.(string), redisPassword.(string), redisDb.(int))

	// 初始化oss
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

	// 启动服务
	app.Run(
		iris.Addr("0.0.0.0:" + strconv.Itoa(port.(int))),
		iris.WithConfiguration(configuration))
}
