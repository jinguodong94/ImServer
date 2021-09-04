package main

import (
	"gindemo/conf"
	"gindemo/dao"
	"gindemo/mq"
	"gindemo/route"
	"gindemo/socket"
	"gindemo/task"
	"github.com/json-iterator/go/extra"
	"log"
)

func main() {

	defer func() {
		//关闭连接
		dao.CloseRedis()
		dao.CloseMysql()
		if mq.MQ != nil {
			mq.MQ.Close()
		}
	}()

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
	go socket.StartServer("localhost:8686")

	//初始化MQ
	mq.MQ = mq.NewRabbitMQ("msg", "amq.topic", conf.Configs.ServerConfig.ServerId, socket.Consumer)
	go mq.MQ.Start()

	//启动http服务
	log.Println("启动http服务")
	route.Route.Run(":8989")
}

func init() {
	extra.RegisterFuzzyDecoders()
	InitLog()
}
