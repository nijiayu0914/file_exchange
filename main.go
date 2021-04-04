// main 项目入口文件，配置加载，数据库初始化等
// 项目基于阿里云OSS封装的文件交换的网盘
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
	// configType 配置文件的来源类型，本地文件存储 or URL GET请求，默认本地文件存储
	configType := flag.String(
		"config_type", // 类型
		"file", // file or http
		"配置文件获取类型，file为本地文件，http访问网络路径")

	// 配置文件存储地址，默认./config/dev/config.yml
	configPath := flag.String(
		"config_path",
		"./config/dev/config.yml", // it can be path or url
		"配置文件地址")
	flag.Parse()

	app := iris.New()
	var configuration iris.Configuration // iris配置对象
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
	port := otherConfig["Port"] // 服务端口
	dbType := otherConfig["DatabaseType"] // 数据库类型
	dbDsn := otherConfig["DatabaseDsn"] // 数据库dsn
	redisDsn := otherConfig["RedisDsn"] // redis dsn
	redisDb := otherConfig["RedisDb"] // 访问redis数据库编号
	redisPassword := otherConfig["RedisPassword"] // redis密码
	oSSEndpoint := otherConfig["OSSEndpoint"] // OSS Endpoint
	oSSAccessKeyID := otherConfig["OSSAccessKeyID"] // OSS AccessKeyID
	oSSAccessKeySecret := otherConfig["OSSAccessKeySecret"] // OSS AccessKeySecret
	oSSBucketName := otherConfig["OSSBucketName"] // OSS BucketName
	oSSRegionId := otherConfig["OSSRegionId"] // OSS RegionId
	oSSRamAccessKeyID := otherConfig["OSSRamAccessKeyID"] // OSS RamAccessKeyID
	oSSRamAccessKeySecret := otherConfig["OSSRamAccessKeySecret"] // OSS RamAccessKeySecret
	oSSRoleArn := otherConfig["OSSRoleArn"] // OSS Role Arn
	oSSRoleSessionName := otherConfig["OSSRoleSessionName"] // OSS RoleSessionName

	// 初始化数据库
	db, err := utils.Db(dbType, dbDsn)
	if err != nil {
		panic("failed to connect database")
	}else{
		db.AutoMigrate(&datamodels.User{}, &datamodels.File{},
		&datamodels.UserPlugin{})
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
