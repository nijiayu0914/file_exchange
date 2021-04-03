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

func (upc *UserPluginController) Read (ctx iris.Context){
	userName := ctx.GetHeader("User-Name")
	userPlugin, err := upc.UserPluginService.FindByUserName(userName)
	if err != nil{
		res := utils.Response{Code: iris.StatusBadRequest,
			Message: "查询失败", Data: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	} else{
		res := utils.Response{Code: iris.StatusOK, Message: "ok", Data: &userPlugin}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}
