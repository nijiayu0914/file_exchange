package controllers

import (
	"file_exchange/datamodels"
	"file_exchange/services"
	"file_exchange/utils"
	"github.com/kataras/iris/v12"
)

// FileController 文件操作控制器
type FileController struct {
	FileService services.IFileService // file服务接口
}

// CheckUuid 检查uuid，用户与用户之间进行隔离，防止非法访问
func (fc *FileController) CheckUuid(ctx iris.Context,
	fileService services.IFileService){
	var rquuid = utils.RequestCheckUuid{}
	userName := ctx.GetHeader("User-Name")
	uuid := ctx.URLParam("uuid")
	ctx.ReadJSON(&rquuid)
	if uuid != ""{
		file, err :=fileService.FindFileByUuid(uuid)
		if err != nil{
			res := utils.Response{Code: iris.StatusForbidden,
				Message: "查询失败"}
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(res)
			return
		}
		if file.UserName != userName{
			res := utils.Response{Code: iris.StatusForbidden,
				Message: "用户匹配错误"}
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(res)
			return
		}
	}
	if rquuid.FileUuid != ""{
		file, err :=fileService.FindFileByUuid(rquuid.FileUuid)
		if err != nil{
			res := utils.Response{Code: iris.StatusForbidden,
				Message: "查询失败"}
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(res)
			return
		}
		if file.UserName != userName{
			res := utils.Response{Code: iris.StatusForbidden,
				Message: "用户匹配错误"}
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(res)
			return
		}
	}
	ctx.Next()
}

// CreateFile 创建文件夹
func (fc *FileController) CreateFile (ctx iris.Context, capacity float64){
	file:= datamodels.File{}
	ctx.ReadJSON(&file)
	file.UserName = ctx.GetHeader("User-Name")
	file.Capacity = capacity
	fileId, fileUuid, err := fc.FileService.CreateFile(&file)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "创建失败", Data: err.Error()}
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

// FindFilesByUserName 根据用户名查找文件
func (fc *FileController) FindFilesByUserName (ctx iris.Context){
	userName := ctx.GetHeader("User-Name")
	files, err := fc.FileService.FindFilesByUserName(userName)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: &files}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// ChangeFileName 更改文件夹名称
func (fc *FileController) ChangeFileName (ctx iris.Context){
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name")
	err := fc.FileService.UpdateFileName(fileName, fileUuid)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "更改失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK,
			Message: "ok", Data: iris.Map{
			"uuid": fileUuid,
			"file_name": fileName,
		}}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// DeleteByUuid 根据uuid删除文件夹
func (fc *FileController) DeleteByUuid (ctx iris.Context){
	fileUuid := ctx.URLParam("uuid")
	err := fc.FileService.DeleteByUuid(fileUuid)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "数据库删除失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "删除成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// UpdateUsage 更新用量
func (fc *FileController) UpdateUsage (ctx iris.Context){
	rquu := utils.RequestUpdateUsage{}
	ctx.ReadJSON(&rquu)
	err := fc.FileService.UpdateUsageCapacity(rquu.UsageCapacity,
		rquu.FileUuid, rquu.How)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "更新失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "更新成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// UpdateCapacity 更新用量上限
func (fc *FileController) UpdateCapacity (ctx iris.Context){
	rqcap := utils.RequestCapacity{}
	ctx.ReadJSON(&rqcap)
	err := fc.FileService.UpdateCapacity(rqcap.Capacity, rqcap.FileUuid)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "更新失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "更新成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// CheckCapacity 检查用量
func (fc *FileController) CheckCapacity (ctx iris.Context){
	fileUuid := ctx.URLParam("uuid")
	usageCapacity, capacity, free, err := fc.FileService.CheckCapacity(fileUuid)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "查询成功",
			Data: iris.Map{
			"usageCapacity": usageCapacity,
			"capacity": capacity,
			"free": free,
		}}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}
