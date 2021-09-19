package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"imserver/constant"
	"imserver/dao"
	"imserver/req"
	"imserver/response"
	"imserver/utils"
	"strconv"
	"time"
)

type UserController struct {
	BaseController
}

//登录
func (userController UserController) Login(ctx *gin.Context) {
	loginReq := &req.LoginReq{}
	err := ctx.ShouldBind(loginReq)
	if err != nil {
		responseOk(ctx, response.NewErrorResponse("参数有误"))
		return
	}
	user := &dao.Users{}
	result := dao.Db.Model(user).Where("account = ? and pwd = ? and deleted_at is null", loginReq.Account, loginReq.Pwd).First(user)

	if result.RowsAffected <= 0 {
		responseOk(ctx, response.NewErrorResponse("登录失败，账户或者密码错误"))
		return
	}

	//生成token
	token := utils.TokenUtils.CreateToken(user.ID)

	dao.Rdb.Set(context.Background(),
		fmt.Sprintf("%s:%s", constant.Redis_key_user_login_token, strconv.FormatUint(uint64(user.ID), 10)),
		token, time.Hour*24*30)

	userInfo := &response.UserInfo{
		Account:  user.Account,
		NickName: user.NickName,
		Icon:     user.Icon,
		Token:    token,
	}
	responseOk(ctx, response.NewSuccessResponse("登录成功", userInfo))
}

//注册
func (userController UserController) Register(ctx *gin.Context) {
	registerReq := &req.RegisterReq{}
	err := ctx.ShouldBind(registerReq)
	if err != nil {
		responseOk(ctx, response.NewErrorResponse("参数有误"))
		return
	}

	users := &dao.Users{
		Account:  registerReq.Account,
		Pwd:      registerReq.Pwd,
		NickName: registerReq.NickName,
		Icon:     registerReq.Icon,
	}
	dao.Db.AutoMigrate(users)

	result := dao.Db.Create(users)
	if result.Error != nil {
		//用户名存在
		responseOk(ctx, response.NewErrorResponse("用户名已存在"))
		return
	}

	//生成token
	token := utils.TokenUtils.CreateToken(users.ID)
	dao.Rdb.Set(context.Background(),
		fmt.Sprintf("%s:%s", constant.Redis_key_user_login_token, strconv.FormatUint(uint64(users.ID), 10)),
		token, time.Hour*24*30)

	userInfo := &response.UserInfo{
		Account:  users.Account,
		NickName: users.NickName,
		Icon:     users.Icon,
		Token:    token,
	}

	responseOk(ctx, response.NewSuccessResponse("注册成功", userInfo))
}

//获取个人信息
func (userController UserController) GetUserInfo(ctx *gin.Context) {
	userId, ok := ctx.GetQuery("userId")
	if !ok {
		responseOk(ctx, response.NewErrorResponse("参数有误"))
		return
	}
	userInfo := &dao.Users{}
	tx := dao.Db.Model(userInfo).Where("id = ?", userId).First(userInfo)
	if tx.Error != nil {
		responseOk(ctx, response.NewErrorResponse("获取失败"))
		return
	}
	userInfo.Pwd = ""
	responseOk(ctx, response.NewSuccessResponse("获取成功", userInfo))
}

//修改个人信息
func (userController UserController) UpdateUserInfo(ctx *gin.Context) {

}
