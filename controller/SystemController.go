package controller

import (
	"gindemo/socket"
	"github.com/gin-gonic/gin"
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
