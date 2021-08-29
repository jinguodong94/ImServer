package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct {
}

func responseOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}
