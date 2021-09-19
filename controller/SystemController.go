package controller

import (
	"github.com/gin-gonic/gin"
	"imserver/socket"
	"net/http"
)

type SystemController struct {
	BaseController
}

func (SystemController) GetSystemStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"allClient":     len(socket.ClientMgr.GetAllClient()),
		"loginedClient": len(socket.ClientMgr.GetLoginClient()),
	})
}
