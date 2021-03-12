package controllers

import (
	"file_exchange/services"
	"file_exchange/utils"
	"github.com/kataras/iris/v12"
	"log"
	"strings"
)

func CreateTestFile(ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name")
	content := ctx.URLParam("content")
	err := ossOperator.CreateStringFile(fileUuid, fileName, content, "file")
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "创建失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: "创建成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func CreateFolder (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name") + "/"
	err := ossOperator.CreateStringFile(fileUuid, fileName, "", "folder")
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "创建失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: "创建成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func ListFiles (ctx iris.Context, ossOperator *services.OssOperator) {
	rqlf := utils.RequestListFiles{}
	ctx.ReadJSON(&rqlf)
	objectsContainer, dirsContainer, err := ossOperator.ListFiles(
		rqlf.FileUuid, rqlf.Path, rqlf.Delimiter)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "读取失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		rplf := utils.ResponseListFiles{
			FilesCount: len(objectsContainer),
			DirsCount: len(dirsContainer),
			Files: objectsContainer,
			Dirs: dirsContainer,
		}
		res := utils.Response{Code: iris.StatusOK,
			Message: "读取成功", Data: &rplf}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func IsFileExist (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name")
	isExist, err := ossOperator.IsExist(fileUuid, fileName)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "判断失败", Data: isExist}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "判断成功", Data: isExist}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func RenameObject (ctx iris.Context, ossOperator *services.OssOperator) {
	rqrename := utils.RequestRenameObject{}
	ctx.ReadJSON(&rqrename)
	err := ossOperator.RenameObject(rqrename.ObjectName, rqrename.NewName)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "修改", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "修改成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func DeleteFile (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name")
	err := ossOperator.DeleteFile(fileUuid, fileName)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "删除失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "删除成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func DeleteChildFile (ctx iris.Context, ossOperator *services.OssOperator) {
	rqlf := utils.RequestListFiles{}
	ctx.ReadJSON(&rqlf)
	err := ossOperator.DeleteChildFile(rqlf.FileUuid, rqlf.Path)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "删除失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "删除成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func DeleteFiles (ctx iris.Context, ossOperator *services.OssOperator) {
	rqdfs := utils.RequestDeleteFiles{}
	ctx.ReadJSON(&rqdfs)
	deleteMarket, err := ossOperator.DeleteFiles(rqdfs.FileUuid, rqdfs.FileNames)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "删除失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK,
			Message: "删除成功", Data: deleteMarket}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func ListDeleteMarkers (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	delimiter := ctx.URLParam("delimiter")
	markers, err := ossOperator.ListDeleteMarkers(fileUuid, "", delimiter)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK,
			Message: "查询成功", Data: markers}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func ListFileVersion (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	path := ctx.URLParam("path")
	objects, err := ossOperator.ListFileVersion(fileUuid, path)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK,
			Message: "查询成功", Data: objects}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func DeleteFileForever (ctx iris.Context, ossOperator *services.OssOperator,
	fileService services.IFileService) {
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name")
	size, err := ossOperator.DeleteFileForever(fileUuid, fileName)
	go func() {
		err := fileService.UpdateUsageCapacity(size, fileUuid, "increase")
		if err != nil{
			log.Println("更新容量失败")
		}
	}()
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "删除失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "删除成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func DeleteFilesForever (ctx iris.Context, ossOperator *services.OssOperator,
	fileService services.IFileService) {
	rqdfs := utils.RequestDeleteFiles{}
	ctx.ReadJSON(&rqdfs)
	size, err := ossOperator.DeleteFilesForever(rqdfs.FileUuid, rqdfs.FileNames)
	go func() {
		err := fileService.UpdateUsageCapacity(size, rqdfs.FileUuid, "increase")
		if err != nil{
			log.Println("更新容量失败")
		}
	}()
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "删除失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "删除成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func DeleteHistoryFile (ctx iris.Context, ossOperator *services.OssOperator){
	rdhf := utils.RequestDeleteHistoryFile{}
	ctx.ReadJSON(&rdhf)
	err := ossOperator.DeleteHistoryFile(rdhf.FileUuid, rdhf.Path, rdhf.VersionId)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "删除失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "删除成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func DeleteLibraryForever (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	objectsContainer, _, err := ossOperator.ListFiles(fileUuid, "", "")
	if err != nil {
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询文件失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}
	for _, object := range objectsContainer{
		fileName := strings.Split(object.Basic.Key, fileUuid + "/")[1]
		if fileName == ""{
			continue
		}
		_, err := ossOperator.DeleteFileForever(fileUuid, fileName)
		if err != nil{
			res := utils.Response{Code: iris.StatusBadRequest,
				Message: "删除失败", Data: err.Error()}
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(res)
			return
		}
	}
	_, err = ossOperator.DeleteFileForever(fileUuid, "") //清理下根节点
	if err != nil{
		log.Println(err)
	}
	ctx.Next()
}

func RestoreFile (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	path := ctx.URLParam("path")
	err := ossOperator.RestoreFile(fileUuid, path)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "恢复失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "恢复成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func CopyFile (ctx iris.Context, ossOperator *services.OssOperator,
	fileService services.IFileService) {
	rqcp := utils.RequestCopy{}
	ctx.ReadJSON(&rqcp)
	size, err := ossOperator.Copy(rqcp.OriginFile, rqcp.DestFile, rqcp.VersionId)
	go func() {
		err := fileService.UpdateUsageCapacity(
			size, strings.Split(rqcp.OriginFile, "/")[0], "decrease")
		if err != nil{
			log.Println("更新容量失败")
		}
	}()
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "复制失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "复制成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func MultipleCopy (ctx iris.Context, ossOperator *services.OssOperator,
	fileService services.IFileService) {
	rqmc := utils.RequestMultipleCopy{}
	ctx.ReadJSON(&rqmc)
	failure, size, err := ossOperator.MultipleCopy(rqmc.CopyList)
	go func() {
		err := fileService.UpdateUsageCapacity(size, strings.Split(
			rqmc.CopyList[0].OriginFile, "/")[0], "decrease")
		if err != nil{
			log.Println("更新容量失败")
		}
	}()
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "复制失败", Data: failure}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "复制成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func ReadFilesSize (ctx iris.Context, ossOperator *services.OssOperator){
	rqrfs := utils.RequestReadFilesSize{}
	ctx.ReadJSON(&rqrfs)
	size, err := ossOperator.ReadFilesCapacity(rqrfs.FileUuid, rqrfs.Files)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "查询成功",
			Data: iris.Map{"size": size}}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func ReadAllFilesCapacity (ctx iris.Context, ossOperator *services.OssOperator){
	fileUuid := ctx.URLParam("uuid")
	size, err := ossOperator.ReadAllFilesCapacity(fileUuid)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "查询成功",
			Data: iris.Map{"size": size}}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func CreateSTS (ctx iris.Context, ossOperator *services.OssOperator) {
	response, err := ossOperator.CreateSTS ()
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "授权失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK,
			Message: "授权成功", Data: &response}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

func DownloadUrl (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name")
	url, err := ossOperator.DownloadFile(fileUuid, fileName)
	if err != nil{
		res := utils.Response{Code: iris.StatusForbidden,
			Message: "授权失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "授权成功", Data: url}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}
