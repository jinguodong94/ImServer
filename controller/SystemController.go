package controller

import (
	"gindemo/ws"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SystemController struct {
	BaseController
}

func (SystemController) GetSystemStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"allClient":     len(ws.ClientMgr.GetAllClient()),
		"loginedClient": len(ws.ClientMgr.GetLoginClient()),
	})
}
