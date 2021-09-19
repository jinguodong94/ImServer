package controller

import (
	"github.com/gin-gonic/gin"
	"imserver/constant"
	"imserver/utils"
	"net/http"
)

type BaseController struct {
}

func responseOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func GetUserIdFromToken(ctx *gin.Context) uint {
	token := ctx.GetHeader(constant.Token)
	id, _ := utils.TokenUtils.GetUserId(token)
	return id
}
