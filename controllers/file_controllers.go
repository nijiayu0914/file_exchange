package controllers

import (
	"file_exchange/datamodels"
	"file_exchange/services"
	"file_exchange/utils"
	"github.com/go-redis/redis"
	"github.com/kataras/iris/v12"
)

type FileController struct {
	FileService services.IFileService
}

func (fc *FileController) CheckToken(ctx iris.Context, redisClient *redis.Client){
	token := ctx.GetHeader("Authorization")
	userName := ctx.GetHeader("User-Name")
	tokenCache, err := redisClient.Get(userName).Result()
	if err != nil{
		res := utils.Response{Code: iris.StatusForbidden, Message: "token验证失败"}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
		return
	}

	checked, _ := utils.Verification(userName, tokenCache)
	if !checked || token != tokenCache{
		res := utils.Response{Code: iris.StatusForbidden, Message: "token验证失败"}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
		return
	}
	ctx.Next()
}

func (fc *FileController) CreateFile (ctx iris.Context){
	file:= datamodels.File{}
	ctx.ReadJSON(&file)
	file.UserName = ctx.GetHeader("User-Name")
	fileId, fileUuid, err := fc.FileService.CreateFile(&file)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest, Message: "创建失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		file.ID = fileId
		file.Uuid = fileUuid
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: &file}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func (fc *FileController) FindFilesByUserName (ctx iris.Context){
	userName := ctx.GetHeader("User-Name")
	files, err := fc.FileService.FindFilesByUserName(userName)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest, Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: &files}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func (fc *FileController) ChangeFileName (ctx iris.Context){
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name")
	err := fc.FileService.UpdateFileName(fileName, fileUuid)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest, Message: "更改失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: iris.Map{
			"uuid": fileUuid,
			"file_name": fileName,
		}}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func (fc *FileController) DeleteByUuid (ctx iris.Context){
	fileUuid := ctx.URLParam("uuid")
	err := fc.FileService.DeleteByUuid(fileUuid)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest, Message: "数据库删除失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "删除成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}
