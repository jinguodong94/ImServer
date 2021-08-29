package main

import (
	"gindemo/conf"
	"gindemo/dao"
	"gindemo/mq"
	"gindemo/plog"
	"gindemo/route"
	"gindemo/task"
	"gindemo/ws"
	"log"
)

func main() {

	defer func() {
		//关闭连接
		dao.CloseRedis()
		dao.CloseMysql()
		mq.Close()
	}()

	plog.Init()

	//初始化配置
	conf.Init()

	//初始化redis
	dao.InitRedis()

	//初始化数据库连接
	dao.InitMysql()

	//初始化路由
	route.Init()

	//启动定时清理客户端任务
	task.StartClearClientTask()

	//启动websocket 服务
	go ws.StartServer("localhost:8686")

	mq.Init()

	//启动http服务
	log.Println("启动http服务")
	route.Route.Run(":8989")
}
