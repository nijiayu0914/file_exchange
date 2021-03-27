// controllers 服务控制层
package controllers

import (
	"file_exchange/datamodels"
	"file_exchange/services"
	"file_exchange/utils"
	"github.com/go-redis/redis"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// UserController 用户操作控制器
type UserController struct {
	UserService services.IUserService // user服务接口
}

// CheckToken 检查用户token信息
func (uc *UserController) CheckToken(ctx iris.Context,
	redisClient *redis.Client){
	token := ctx.GetHeader("Authorization")
	userName := ctx.GetHeader("User-Name")
	tokenCache, err := redisClient.Get(userName).Result()
	if err != nil{
		res := utils.Response{Code: iris.StatusForbidden,
			Message: "token验证失败"}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
		return
	}

	checked, _ := utils.Verification(userName, tokenCache)
	if !checked || token != tokenCache{
		res := utils.Response{Code: iris.StatusForbidden,
			Message: "token验证失败"}
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(res)
		return
	}
	ctx.Next()
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(ctx iris.Context,
	redisClient *redis.Client){
	user := datamodels.User{}
	ctx.ReadJSON(&user)
	userId, err := uc.UserService.AddUser(&user)
	var res utils.Response
	if err != nil{
		res = utils.Response{Code: iris.StatusBadRequest,
			Message: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
	}else{
		token, err := utils.GenerateToken(user.UserName, user.Password)
		if err != nil{
			res = utils.Response{Code: iris.StatusBadRequest,
				Message: err.Error()}
			ctx.StatusCode(iris.StatusBadRequest)
		}else{
			err := redisClient.Set(user.UserName, token, 48 * time.Hour).Err()
			if err != nil{
				res = utils.Response{Code: iris.StatusBadRequest,
					Message: "缓存失败", Data: err.Error()}
				ctx.StatusCode(iris.StatusBadRequest)
			}else{
				res = utils.Response{Code: iris.StatusOK,
					Message: "ok", Data: iris.Map{
					"user_id": userId, "token": token}}
				ctx.StatusCode(iris.StatusOK)
			}
		}
	}
	ctx.JSON(res)
}

// Login 用户登录
func (uc *UserController) Login(ctx iris.Context, redisClient *redis.Client,
	admin string, adminPassword string){
	var res utils.Response
	user := datamodels.User{}
	ctx.ReadJSON(&user)
	if user.UserName == admin{
		if user.Password != utils.Str2md5(adminPassword){
			ctx.JSON(utils.Response{Code: iris.StatusForbidden,
				Message: "管理员密码错误"})
			return
		}else{
			token, err := utils.GenerateToken(admin, adminPassword)
			if err != nil{
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(utils.Response{Code: iris.StatusBadRequest,
					Message: err.Error()})
				return
			}else{
				err := redisClient.Set(user.UserName,
					token, 48 * time.Hour).Err()
				if err != nil{
					ctx.StatusCode(iris.StatusBadRequest)
					ctx.JSON(utils.Response{Code: iris.StatusBadRequest,
						Message: "缓存失败", Data: err.Error()})
					return
				}else{
					ctx.StatusCode(iris.StatusOK)
					ctx.JSON(utils.Response{Code: iris.StatusOK,
						Message: "ok", Data: iris.Map{"token": token}})
					return
				}
			}
		}
	}
	userRes, err := uc.UserService.FindUser(user.UserName)
	if err != nil{
		res = utils.Response{Code: iris.StatusBadRequest,
			Message: err.Error()}
		ctx.StatusCode(iris.StatusBadRequest)
	}else{
		err := bcrypt.CompareHashAndPassword([]byte(userRes.Password),
			[]byte(user.Password))
		if err != nil{
			res = utils.Response{Code: iris.StatusUnauthorized,
				Message: "密码错误"}
			ctx.StatusCode(iris.StatusUnauthorized)
		}else{
			tokenCache, err := redisClient.Get(user.UserName).Result()
			if err != nil{
				token, err := utils.GenerateToken(
					userRes.UserName, userRes.Password)
				if err != nil{
					res = utils.Response{Code: iris.StatusBadRequest,
						Message: err.Error()}
					ctx.StatusCode(iris.StatusBadRequest)
				}else{
					err := redisClient.Set(user.UserName,
						token, 48 * time.Hour).Err()
					if err != nil{
						res = utils.Response{Code: iris.StatusBadRequest,
							Message: "缓存失败", Data: err.Error()}
						ctx.StatusCode(iris.StatusBadRequest)
					}else{
						res = utils.Response{Code: iris.StatusOK,
							Message: "ok", Data: iris.Map{"token": token}}
						ctx.StatusCode(iris.StatusOK)
					}
				}
			}else{
				checked, _ := utils.Verification(user.UserName, tokenCache)
				if checked{
					res = utils.Response{Code: iris.StatusOK, Message: "ok",
						Data: iris.Map{"token": tokenCache}}
					ctx.StatusCode(iris.StatusOK)
				}else{
					token, err := utils.GenerateToken(userRes.UserName,
						userRes.Password)
					if err != nil{
						res = utils.Response{Code: iris.StatusBadRequest,
							Message: err.Error()}
						ctx.StatusCode(iris.StatusBadRequest)
					}else{
						err := redisClient.Set(user.UserName,
							token, 48 * time.Hour).Err()
						if err != nil{
							res = utils.Response{Code: iris.StatusBadRequest,
								Message: "缓存失败", Data: err.Error()}
							ctx.StatusCode(iris.StatusBadRequest)
						}else{
							res = utils.Response{Code: iris.StatusOK,
								Message: "ok",
								Data: iris.Map{"token": token}}
							ctx.StatusCode(iris.StatusOK)
						}
					}
				}
			}
		}

	}
	ctx.JSON(res)
}

// FindAllUser 查询所有用户
func (uc *UserController) FindAllUser(ctx iris.Context){
	users, err := uc.UserService.FindAllUser()
	if err != nil{
 		res := utils.Response{Code: iris.StatusBadRequest, Message: "获取失败"}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(res)
	}else{
		res := utils.Response{Code: iris.StatusOK, Message: "获取成功",
			Data: iris.Map{"users": users}}
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(res)
	}
}

// ChangePassword 修改密码
func (uc *UserController) ChangePassword(ctx iris.Context){
	var changePassword utils.RequestChangePassword
	ctx.ReadJSON(&changePassword)
	userRes, err := uc.UserService.FindUser(changePassword.UserName)
	if err != nil{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(utils.Response{Code: 400, Message: err.Error()})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userRes.Password),
		[]byte(changePassword.OldPassword))
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(utils.Response{Code: iris.StatusUnauthorized,
			Message: "密码错误"})
		return
	}
	userId, err := uc.UserService.ChangePassword(userRes,
		changePassword.NewPassword)
	if err != nil{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(utils.Response{Code: iris.StatusBadRequest,
			Message: "密码修改失败"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(utils.Response{Code: iris.StatusOK, Message: "密码修改成功",
		Data: iris.Map{"user_id": userId}})
}

func (uc *UserController) ResetPassword(ctx iris.Context){
	var changePassword utils.RequestResetPassword
	ctx.ReadJSON(&changePassword)
	userRes, err := uc.UserService.FindUser(changePassword.UserName)
	if err != nil{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(utils.Response{Code: 400, Message: err.Error()})
		return
	}
	userId, err := uc.UserService.ChangePassword(userRes,
		changePassword.NewPassword)
	if err != nil{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(utils.Response{Code: iris.StatusBadRequest,
			Message: "密码修改失败"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(utils.Response{Code: iris.StatusOK, Message: "密码修改成功",
		Data: iris.Map{"user_id": userId}})
}

// DelToken 删除缓存token
func (uc *UserController) DelToken(ctx iris.Context,
	redisClient *redis.Client){
	userName := ctx.GetHeader("User-Name")
	_, err := redisClient.Del(userName).Result()
	if err != nil{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(utils.Response{Code: iris.StatusBadRequest,
			Message: "删除失败"})
	}else{
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(utils.Response{Code: iris.StatusOK, Message: "删除成功"})
	}
}
