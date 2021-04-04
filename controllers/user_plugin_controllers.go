package controllers

import (
	"file_exchange/services"
	"file_exchange/utils"
	"github.com/kataras/iris/v12"
)

// UserPluginController 用户配置操作控制器
type UserPluginController struct {
	UserPluginService services.IUserPluginService // 用户配置服务接口
}

// checkPermission 后台数据访问鉴权
func (upc *UserPluginController) CheckPermission(
	ctx iris.Context, admin string){
	userName := ctx.GetHeader("User-Name")
	if userName == admin{
		ctx.Next()
	}
	userPlugin, err := upc.UserPluginService.FindByUserName(userName)
	if err != nil {
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "鉴权查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
		return
	}
	if userPlugin.Permission == 1004{
		ctx.Next()
	}else{
		res := utils.Response{Code: iris.StatusForbidden,
			Message: "无权限", Data: "无权限"}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
		return
	}
}

// Read 读取单个用户配置信息
func (upc *UserPluginController) Read (ctx iris.Context, admin string){
	userName := ctx.GetHeader("User-Name")
	userPlugin, err := upc.UserPluginService.FindByUserName(userName)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: iris.Map{
			"user_plugin": &userPlugin,
			"admin_name": admin,
		}}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// ReadAll 分页读取所有用户配置信息
func (upc *UserPluginController) ReadAll (ctx iris.Context){
	page, err := ctx.URLParamInt("page")
	if err != nil{
		page = 1
	}
	pageSize, err := ctx.URLParamInt("page_size")
	if err != nil{
		pageSize = 1
	}
	keyWord := ctx.URLParam("key_word")
	userPlugins, errPlugin := upc.UserPluginService.FindByPaginate(
		page, pageSize, keyWord)
	count, errCount := upc.UserPluginService.Count(keyWord)
	if errPlugin != nil || errCount != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: nil}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: iris.Map{
			"count": &count,
			"page": &page,
			"page_size": &pageSize,
			"key_word": &keyWord,
			"user_plugins": &userPlugins,
		}}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// UpdatePermission 更新用户权限
func (upc *UserPluginController) UpdatePermission (ctx iris.Context){
	rqupm := utils.RequestUpdatePermission{}
	ctx.ReadJSON(&rqupm)
	err := upc.UserPluginService.UpdatePermission(rqupm.UserName, rqupm.Permission)
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

// UpdateMaxLibrary 更新用户Library最大数量
func (upc *UserPluginController) UpdateMaxLibrary (ctx iris.Context){
	rquml := utils.RequestUpdateMaxLibrary{}
	ctx.ReadJSON(&rquml)
	err := upc.UserPluginService.UpdateMaxLibrary(rquml.UserName, rquml.MaxLibrary)
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
