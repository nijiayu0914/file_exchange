// routes 服务的路由定义， 所有API接口均定义在此处
package routes

import (
	"file_exchange/controllers"
	"file_exchange/repositories"
	"file_exchange/services"
	"github.com/go-redis/redis"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

// RegisterUserCollercors 注册用户控制器
func RegisterUserCollercors(db *gorm.DB) services.IUserService{
	repository := repositories.NewUserRepository(db)
	service := services.NewUserService(repository)
	return service
}

// RegisterFileCollercors 注册文件控制器
func RegisterFileCollercors(db *gorm.DB) services.IFileService{
	repository := repositories.NewFileRepository(db)
	service := services.NewFileService(repository)
	return service
}

// Routes 路由接口，所有API由此函数生成
func Routes(
	app *iris.Application, // iris app
	db *gorm.DB, // gorm Db
	redisClient *redis.Client, // redis client
	ossOperator *services.OssOperator, // OSS操作对象
	otherConfig map[string]interface{}, // 自定义配置
	){
	admin := otherConfig["Admin"].(string) // 超级管理员用户名
	adminPassword := otherConfig["AdminPassword"].(string) // 超级管理员密码
	userService := RegisterUserCollercors(db) // 用户服务
	fileService := RegisterFileCollercors(db) // 文件服务

	// 服务状态健康检查
	// return:
	//	message: 状态正常
	app.Get("/health", func(ctx iris.Context){
		ctx.JSON(iris.Map{"message": "状态正常"})
	})

	// 用户相关API根路由
	users := app.Party("/user")
	// 获取所有用户信息
	// Header:
	//	Authorization: token
	//  User-Name: 用户名
	// return:
	//	message: 状态正常
	users.Get("/users", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.FindAllUser(ctx)
	})
	// 删除token缓存
	// Header:
	//  User-Name: 用户名
	users.Delete("/deltoken", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.DelToken(ctx, redisClient)
	})
	// 创建用户
	// Request Body:
	//  user_name: 用户名
	//  password: 密码
	users.Post("/create", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.CreateUser(ctx, redisClient)
	})
	// 用户登录
	// Request Body:
	//  user_name: 用户名
	//  password: 密码
	users.Post("/login", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.Login(ctx, redisClient, admin, adminPassword)
	})
	// 修改密码
	// Request Body:
	//  user_name: 用户名
	//  old_password: 旧密码
	//  new_password: 新密码
	users.Post("/change_password", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.ChangePassword(ctx)
	})
	// 重置密码
	// Request Body:
	//  user_name: 用户名
	//  new_password: 新密码
	users.Post("/reset_password", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.ResetPassword(ctx)
	})

	// 文件相关API根路由
	files := app.Party("/file",
		// token检查
		// Header:
		//	Authorization: token
		//  User-Name: 用户名
		func(ctx iris.Context){
			userController := controllers.UserController{
				UserService: userService}
			userController.CheckToken(ctx, redisClient)
		},
		// uuid检查
		// 会查找url params中uuid的值和request body中file_uuid的值
		// 并与用户名做匹配，判断此uuid对应的文件夹是否属于该用户
		func(ctx iris.Context){
			fileController := controllers.FileController{
				FileService: fileService}
			fileController.CheckUuid(ctx, fileService)
		})
	// 更改根文件夹（用户级）名称，并分配uuid
	// Params:
	//	uuid: 文件夹uuid
	//  file_name: 新文件名称
	files.Get("/change_library_name", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.ChangeFileName(ctx)
	})
	// 根据用户名查询所有根文件夹（用户级）
	// Header:
	//  User-Name: 用户名
	files.Get("/find_libraries/byuserName", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.FindFilesByUserName(ctx)
	})
	// 删除根文件夹（用户级）
	// Params:
	//	uuid: 文件夹uuid
	files.Get("/delete_library", func(ctx iris.Context){
		controllers.DeleteLibraryForever(ctx, ossOperator)
	}, func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.DeleteByUuid(ctx)
	})
	// 创建测试文件
	// Params:
	//	uuid: 文件夹uuid
	//  file_name: 文件名称
	//  content: 文件内容
	files.Get("/create_test_file", func(ctx iris.Context){
		controllers.CreateTestFile(ctx, ossOperator)
	})
	// 创建文件夹
	// OSS没有文件夹概念，所有文件层级都是一致的，
	// 这里创建一个空文件，文件名会以"/"结尾，用于模拟文件夹
	// Params:
	//	uuid: 文件夹uuid
	//  file_name: 文件名称
	files.Get("/create_folder", func(ctx iris.Context){
		controllers.CreateFolder(ctx, ossOperator)
	})
	// 判断文件是否存在
	// Params:
	//	uuid: 文件夹uuid
	//  file_name: 文件名称
	files.Get("/file_exist", func(ctx iris.Context){
		controllers.IsFileExist(ctx, ossOperator)
	})
	// 列举删除标记
	// Params:
	//	uuid: 文件夹uuid
	//  delimiter: 查询终止符，可添加"/"仅查询当前文件夹下的删除标记
	files.Get("/list_delete_markers", func(ctx iris.Context){
		controllers.ListDeleteMarkers(ctx, ossOperator)
	})
	// 列举文件版本
	// Params:
	//	uuid: 文件夹uuid
	//  path: 文件路径
	files.Get("/list_file_version", func(ctx iris.Context){
		controllers.ListFileVersion(ctx, ossOperator)
	})
	// STS临时授权
	files.Get("/grant", func(ctx iris.Context){
		controllers.CreateSTS(ctx, ossOperator)
	})
	// 读取所有文件大小
	// Params:
	//	uuid: 文件夹uuid
	files.Get("/read_allfiles_size", func(ctx iris.Context){
		controllers.ReadAllFilesCapacity(ctx, ossOperator)
	})
	// 检查当前用量
	// Params:
	//	uuid: 文件夹uuid
	files.Get("/check_capacity", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.CheckCapacity(ctx)
	})
	// 返回临时的文件下载url
	// Params:
	//	uuid: 文件夹uuid
	files.Get("/download", func(ctx iris.Context){
		controllers.DownloadUrl(ctx, ossOperator)
	})
	// 永久删除文件
	// Params:
	//	uuid: 文件夹uuid
	//  file_name: 文件名称
	files.Delete("/delete_file_forever", func(ctx iris.Context){
		controllers.DeleteFileForever(ctx, ossOperator, fileService)
	})
	// 更新用量
	// Request Body:
	//	file_uuid: 文件夹uuid
	//  usage_capacity: 更新的用量
	//  how: increase, decrease, overwrite
	files.Put("/update_usage", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.UpdateUsage(ctx)
	})
	// 更新容量
	// Request Body:
	//	file_uuid: 文件夹uuid
	//  capacity: 更新的容量
	files.Put("/update_capacity", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.UpdateCapacity(ctx)
	})
	// 还原文件
	// Params:
	//	uuid: 文件夹uuid
	//  path: 文件路径
	files.Put("/restore_file", func(ctx iris.Context){
		controllers.RestoreFile(ctx, ossOperator)
	})
	// 重命名文件
	// Request Body:
	//	file_uuid: 文件uuid
	//	object_name: OSS对象文件名
	//  new_name: 需重命名的文件名
	files.Put("/rename_file", func(ctx iris.Context){
		controllers.RenameObject(ctx, ossOperator)
	})
	// 删除文件
	// Params:
	//	uuid: 文件夹uuid
	//  file_name: 文件名称
	files.Delete("/delete_file", func(ctx iris.Context){
		controllers.DeleteFile(ctx, ossOperator)
	})
	// 创建根文件夹（用户级）
	// Header:
	//  User-Name: 用户名
	// Request Body:
	//  file_name: 文件名
	files.Post("/create_library", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.CreateFile(ctx, otherConfig["UserCapacity"].(float64))
	})
	// 列举文件
	// Request Body:
	//	file_uuid: 文件uuid
	//  path: 子文件路径
	//  delimiter: 终止符
	files.Post("/list_files", func(ctx iris.Context){
		controllers.ListFiles(ctx, ossOperator)
	})
	// 删除多个文件
	// Request Body:
	//	file_uuid: 文件uuid
	//  file_names: 文件名数组
	files.Post("/delete_files", func(ctx iris.Context){
		controllers.DeleteFiles(ctx, ossOperator)
	})
	// 永久删除多个文件
	// Request Body:
	//	file_uuid: 文件uuid
	//  file_names: 文件名数组
	files.Post("/delete_files_forever", func(ctx iris.Context){
		controllers.DeleteFilesForever(ctx, ossOperator, fileService)
	})
	// 删除历史文件
	// Request Body:
	//  file_uuid: 文件uuid
	//  path: 子文件路径
	//  version_id: 版本号
	files.Post("/delete_history_file", func(ctx iris.Context){
		controllers.DeleteHistoryFile(ctx, ossOperator)
	})
	// 删除子文件夹
	// Request Body:
	//	file_uuid: 文件uuid
	//  path: 子文件路径
	//  delimiter: 终止符
	files.Post("/delete_child_file", func(ctx iris.Context){
		controllers.DeleteChildFile(ctx, ossOperator)
	})
	// 拷贝文件
	// Request Body:
	//	file_uuid: 文件uuid
	//	origin_file: 源文件
	//  dest_file: 目标文件
	//  version_id: 源文件版本号
	files.Post("/copy_file", func(ctx iris.Context){
		controllers.CopyFile(ctx, ossOperator, fileService)
	})
	// 拷贝多个文件
	// Request Body:
	//	file_uuid: 文件uuid
	//	copy_list: 拷贝文件对象的数组
	files.Post("/copy_multi_files", func(ctx iris.Context){
		controllers.MultipleCopy(ctx, ossOperator, fileService)
	})
	// 读取多个文件大小总和
	// Request Body:
	//	file_uuid: 文件uuid
	//  files: 读取文件的列表
	//    	file_name: 文件名
	//      version_id: 版本号
	files.Post("/read_files_size", func(ctx iris.Context){
		controllers.ReadFilesSize(ctx, ossOperator)
	})
}
