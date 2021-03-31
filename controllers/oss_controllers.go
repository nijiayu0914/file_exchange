package controllers

import (
	"encoding/json"
	"file_exchange/services"
	"file_exchange/utils"
	"github.com/go-redis/redis"
	"github.com/kataras/iris/v12"
	"log"
	"strings"
	"time"
)

// CreateTestFile 创建测试文件
func CreateTestFile(ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name")
	content := ctx.URLParam("content")
	err := ossOperator.CreateStringFile(fileUuid, fileName,
		content, "file")
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "创建失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok",
			Data: "创建成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// CreateFolder 创建文件夹
func CreateFolder (ctx iris.Context, ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	fileName := ctx.URLParam("file_name") + "/"
	err := ossOperator.CreateStringFile(fileUuid, fileName,
		"", "folder")
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "创建失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK,
			Message: "ok", Data: "创建成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// ListFiles 列举文件
func ListFiles (ctx iris.Context, ossOperator *services.OssOperator,
	redisClient *redis.Client) {
	userName := ctx.GetHeader("User-Name")
	rqlf := utils.RequestListFiles{}
	ctx.ReadJSON(&rqlf)
	var objectsContainer []utils.ObjectInfoCollection
	var dirsContainer []utils.DirInfoCollection
	var err error
	rplf := utils.ResponseListFiles{}
	if rqlf.Force == true{
		objectsContainer, dirsContainer, err = ossOperator.ListFiles(
			rqlf.FileUuid, rqlf.Path, rqlf.Delimiter)
		rplf = utils.ResponseListFiles{
			FilesCount: len(objectsContainer),
			DirsCount: len(dirsContainer),
			Files: objectsContainer,
			Dirs: dirsContainer,
		}
		if err != nil{
			res := utils.Response{Code: iris.StatusBadRequest,
				Message: "读取失败", Data: err.Error()}
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(res)
			return
		}
	}else{
		key := utils.ListFilesRedisKey(rqlf.FileUuid, rqlf.Path,
			rqlf.Delimiter, userName)
		cacheListFiles, err := redisClient.Get(key).Result()
		if err != nil{
			objectsContainer, dirsContainer, err = ossOperator.ListFiles(
				rqlf.FileUuid, rqlf.Path, rqlf.Delimiter)
			rplf = utils.ResponseListFiles{
				FilesCount: len(objectsContainer),
				DirsCount: len(dirsContainer),
				Files: objectsContainer,
				Dirs: dirsContainer,
			}
		}else{
			json.Unmarshal([]byte(cacheListFiles), &rplf)
		}
	}
	go func(){
		key := utils.ListFilesRedisKey(rqlf.FileUuid, rqlf.Path,
			rqlf.Delimiter, userName)
		value, err := json.Marshal(rplf)
		if err != nil{
			log.Println("数据列表解析失败")
		}
		err = redisClient.Set(key, string(value), 240 * time.Hour).Err()
		if err != nil{
			log.Println("缓存数据列表失败")
		}
	}()
	res := utils.Response{Code: iris.StatusOK,
			Message: "读取成功", Data: &rplf}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(res)
	}


// IsFileExist 检查文件是否存在
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
		res := utils.Response{Code: iris.StatusOK,
			Message: "判断成功", Data: isExist}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// RenameObject 重命名对象
func RenameObject (ctx iris.Context, ossOperator *services.OssOperator) {
	rqrename := utils.RequestRenameObject{}
	ctx.ReadJSON(&rqrename)
	if len(rqrename.FileUuid) == 0 {
		res := utils.Response{Code: iris.StatusForbidden,
			Message: "uuid不能为空"}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
		return
	}
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

// DeleteFile 删除文件
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

// DeleteChildFile 删除文件夹下文件夹
func DeleteChildFile (ctx iris.Context, ossOperator *services.OssOperator) {
	rdf := utils.RequestDeleteChildFile{}
	ctx.ReadJSON(&rdf)
	err := ossOperator.DeleteChildFile(rdf.FileUuid, rdf.Path)
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

// DeleteFiles 多个文件打上删除标记
func DeleteFiles (ctx iris.Context, ossOperator *services.OssOperator) {
	rqdfs := utils.RequestDeleteFiles{}
	ctx.ReadJSON(&rqdfs)
	deleteMarket, err := ossOperator.DeleteFiles(rqdfs.FileUuid,
		rqdfs.FileNames)
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

// ListDeleteMarkers 列举删除标记
func ListDeleteMarkers (ctx iris.Context, ossOperator *services.OssOperator,
	redisClient *redis.Client) {
	userName := ctx.GetHeader("User-Name")
	fileUuid := ctx.URLParam("uuid")
	delimiter := ctx.URLParam("delimiter")
	force, err := ctx.URLParamBool("force")
	var markers []map[string]interface{}
	if err != nil{
		force = false
	}
	if force == true{
		markers, err = ossOperator.ListDeleteMarkers(fileUuid, "", delimiter)
		if err != nil{
			res := utils.Response{Code: iris.StatusBadRequest,
				Message: "查询失败", Data: err.Error()}
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(res)
			return
		}
	}else {
		key := utils.ListDeleteMarkersRedisKey(fileUuid, userName)
		cacheDeleteMarkers, err := redisClient.Get(key).Result()
		if err != nil {
			markers, err = ossOperator.ListDeleteMarkers(
				fileUuid, "", delimiter)
			if err != nil {
				res := utils.Response{Code: iris.StatusBadRequest,
					Message: "查询失败", Data: err.Error()}
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(res)
				return
			}
		}
		json.Unmarshal([]byte(cacheDeleteMarkers), &markers)
	}

	go func(){
		key := utils.ListDeleteMarkersRedisKey(fileUuid, userName)
		value, err := json.Marshal(&markers)
		if err != nil{
			log.Println("删除列表解析失败")
		}
		err = redisClient.Set(key, string(value), 240 * time.Hour).Err()
		if err != nil{
			log.Println("缓存删除列表失败")
		}
	}()
	res := utils.Response{Code: iris.StatusOK,
		Message: "查询成功", Data: &markers}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(res)

}

// ListFileVersion 列举文件版本
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

// DeleteFileForever 永久删除文件
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

// DeleteFilesForever 永久删除多个文件
func DeleteFilesForever (ctx iris.Context, ossOperator *services.OssOperator,
	fileService services.IFileService) {
	rqdfs := utils.RequestDeleteFiles{}
	ctx.ReadJSON(&rqdfs)
	size, err := ossOperator.DeleteFilesForever(rqdfs.FileUuid, rqdfs.FileNames)
	go func() {
		err := fileService.UpdateUsageCapacity(size,
			rqdfs.FileUuid, "increase")
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

// DeleteHistoryFile 删除历史文件
func DeleteHistoryFile (ctx iris.Context, ossOperator *services.OssOperator){
	rdhf := utils.RequestDeleteHistoryFile{}
	ctx.ReadJSON(&rdhf)
	err := ossOperator.DeleteHistoryFile(rdhf.FileUuid,
		rdhf.Path, rdhf.VersionId)
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

// DeleteLibraryForever 删除整个用户文件
func DeleteLibraryForever (ctx iris.Context,
	ossOperator *services.OssOperator) {
	fileUuid := ctx.URLParam("uuid")
	objectsContainer, _, err := ossOperator.ListFiles(
		fileUuid, "", "")
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

// RestoreFile 还原文件
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

// CopyFile 复制文件
func CopyFile (ctx iris.Context, ossOperator *services.OssOperator,
	fileService services.IFileService) {
	rqcp := utils.RequestCopy{}
	ctx.ReadJSON(&rqcp)
	if len(rqcp.FileUuid) == 0 {
		res := utils.Response{Code: iris.StatusForbidden,
			Message: "uuid不能为空"}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
		return
	}
	size, err := ossOperator.Copy(rqcp.OriginFile,
		rqcp.DestFile, rqcp.VersionId)
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

// MultipleCopy 复制多个文件
func MultipleCopy (ctx iris.Context, ossOperator *services.OssOperator,
	fileService services.IFileService) {
	rqmc := utils.RequestMultipleCopy{}
	ctx.ReadJSON(&rqmc)
	if len(rqmc.FileUuid) == 0 {
		res := utils.Response{Code: iris.StatusForbidden,
			Message: "uuid不能为空"}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
		return
	}
	failure, size, err := ossOperator.MultipleCopy(rqmc.CopyList)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "复制失败", Data: failure}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		go func() {
			err := fileService.UpdateUsageCapacity(size, strings.Split(
				rqmc.CopyList[0].DestFile, "/")[0], "increase")
			if err != nil{
				log.Println("更新容量失败")
			}
		}()
		res := utils.Response{Code: iris.StatusOK, Message: "复制成功"}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// ReadFilesSize 读取文件容量
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

// ReadAllFilesCapacity 读取所有文件容量
func ReadAllFilesCapacity (ctx iris.Context,
	ossOperator *services.OssOperator){
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

// CreateSTS 创建临时授权
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

// DownloadUrl 获取临时文件下载地址
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
		res := utils.Response{Code: iris.StatusOK,
			Message: "授权成功", Data: url}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}
