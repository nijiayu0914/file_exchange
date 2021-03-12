package routes

import (
	"file_exchange/controllers"
	"file_exchange/repositories"
	"file_exchange/services"
	"github.com/go-redis/redis"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func RegisterUserCollercors(db *gorm.DB) services.IUserService{
	repository := repositories.NewUserRepository(db)
	service := services.NewUserService(repository)
	return service
}

func RegisterFileCollercors(db *gorm.DB) services.IFileService{
	repository := repositories.NewFileRepository(db)
	service := services.NewFileService(repository)
	return service
}

func Routes(app *iris.Application, db *gorm.DB,
	redisClient *redis.Client, ossOperator *services.OssOperator,
	otherConfig map[string]interface{}){
	admin := otherConfig["Admin"].(string)
	adminPassword := otherConfig["AdminPassword"].(string)
	userService := RegisterUserCollercors(db)
	fileService := RegisterFileCollercors(db)

	app.Get("/health", func(ctx iris.Context){
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	users := app.Party("/user")
	users.Get("/users", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.FindAllUser(ctx)
	})
	users.Delete("/deltoken", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.DelToken(ctx, redisClient)
	})
	users.Post("/create", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.CreateUser(ctx, redisClient) //e10adc3949ba59abbe56e057f20f883e 96e79218965eb72c92a549dd5a330112
	})
	users.Post("/login", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.Login(ctx, redisClient, admin, adminPassword)
	})
	users.Post("/change_password", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.ChangePassword(ctx)
	})
	users.Post("/reset_password", func(ctx iris.Context){
		userController := controllers.UserController{UserService: userService}
		userController.ResetPassword(ctx)
	})

	files := app.Party("/file",
		func(ctx iris.Context){
			userController := controllers.UserController{UserService: userService}
			userController.CheckToken(ctx, redisClient)
		},
		func(ctx iris.Context){
			fileController := controllers.FileController{FileService: fileService}
			fileController.CheckUuid(ctx, fileService)
		})
	files.Get("/change_library_name", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.ChangeFileName(ctx)
	})
	files.Get("/find_libraries/byuserName", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.FindFilesByUserName(ctx)
	})
	files.Get("/delete_library", func(ctx iris.Context){
		controllers.DeleteLibraryForever(ctx, ossOperator)
	}, func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.DeleteByUuid(ctx)
	})
	files.Get("/create_test_file", func(ctx iris.Context){
		controllers.CreateTestFile(ctx, ossOperator)
	})
	files.Get("/create_folder", func(ctx iris.Context){
		controllers.CreateFolder(ctx, ossOperator)
	})
	files.Get("/file_exist", func(ctx iris.Context){
		controllers.IsFileExist(ctx, ossOperator)
	})
	files.Get("/list_delete_markers", func(ctx iris.Context){
		controllers.ListDeleteMarkers(ctx, ossOperator)
	})
	files.Delete("/delete_file_forever", func(ctx iris.Context){
		controllers.DeleteFileForever(ctx, ossOperator, fileService)
	})
	files.Get("/list_file_version", func(ctx iris.Context){
		controllers.ListFileVersion(ctx, ossOperator)
	})
	files.Get("/grant", func(ctx iris.Context){
		controllers.CreateSTS(ctx, ossOperator)
	})
	files.Get("/read_allfiles_size", func(ctx iris.Context){
		controllers.ReadAllFilesCapacity(ctx, ossOperator)
	})
	files.Get("/check_capacity", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.CheckCapacity(ctx)
	})
	files.Get("/download", func(ctx iris.Context){
		controllers.DownloadUrl(ctx, ossOperator)
	})
	files.Put("/update_usage", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.UpdateUsage(ctx)
	})
	files.Put("/update_capacity", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.UpdateCapacity(ctx)
	})
	files.Put("/restore_file", func(ctx iris.Context){
		controllers.RestoreFile(ctx, ossOperator)
	})
	files.Put("/rename_file", func(ctx iris.Context){
		controllers.RenameObject(ctx, ossOperator)
	})
	files.Delete("/delete_file", func(ctx iris.Context){
		controllers.DeleteFile(ctx, ossOperator)
	})
	files.Post("/create_library", func(ctx iris.Context){
		fileController := controllers.FileController{FileService: fileService}
		fileController.CreateFile(ctx, otherConfig["UserCapacity"].(float64))
	})
	files.Post("/list_files", func(ctx iris.Context){
		controllers.ListFiles(ctx, ossOperator)
	})
	files.Post("/delete_files", func(ctx iris.Context){
		controllers.DeleteFiles(ctx, ossOperator)
	})
	files.Post("/delete_files_forever", func(ctx iris.Context){
		controllers.DeleteFilesForever(ctx, ossOperator, fileService)
	})
	files.Post("/delete_history_file", func(ctx iris.Context){
		controllers.DeleteHistoryFile(ctx, ossOperator)
	})
	files.Post("/delete_child_file", func(ctx iris.Context){
		controllers.DeleteChildFile(ctx, ossOperator)
	})
	files.Post("/copy_file", func(ctx iris.Context){
		controllers.CopyFile(ctx, ossOperator, fileService)
	})
	files.Post("/copy_multi_files", func(ctx iris.Context){
		controllers.MultipleCopy(ctx, ossOperator, fileService)
	})
	files.Post("/read_files_size", func(ctx iris.Context){
		controllers.ReadFilesSize(ctx, ossOperator)
	})
}
