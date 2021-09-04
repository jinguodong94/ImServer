package controller

import (
	"context"
	"fmt"
	"gindemo/constant"
	"gindemo/dao"
	"gindemo/response"
	"gindemo/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	BaseController
}

//登录校验
func (AuthController) LoginAuth(ctx *gin.Context) {
	token := ctx.GetHeader(constant.Token)
	if token == "" {
		responseOk(ctx, response.NewNotLoginErrorResponse())
		ctx.Abort()
		return
	}

	userId, err := utils.TokenUtils.GetUserId(token)
	if err != nil {
		responseOk(ctx, response.NewNotLoginErrorResponse())
		ctx.Abort()
		return
	}

	key := fmt.Sprintf("%s:%d", constant.Redis_key_user_login_token, userId)
	result := dao.Rdb.Get(context.Background(), key)
	t, err := result.Result()
	if err != nil || token != t {
		responseOk(ctx, response.NewNotLoginErrorResponse())
		ctx.Abort()
		return
	}
	ctx.Next()
}
