package main

import (
	"fmt"
	"github.com/json-iterator/go/extra"
	"imserver/conf"
	"imserver/dao"
	"imserver/mq"
	"imserver/route"
	"imserver/socket"
	"imserver/task"
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

	serverConfig := conf.Configs.ServerConfig

	//启动websocket 服务
	go socket.StartServer(fmt.Sprintf("%s:%d", serverConfig.WebSocketServerIp, serverConfig.WebSocketServerPort))

	//初始化MQ
	mq.MQ = mq.NewRabbitMQ("msg", "amq.topic", serverConfig.ServerId, socket.Consumer)
	go mq.MQ.Start()

	//启动http服务
	log.Println("启动http服务")
	err := route.Route.Run(fmt.Sprintf("%s:%d", serverConfig.HttpServerIp, serverConfig.HttpServerPort))
	if err != nil {
		panic("httpServer run error : " + err.Error())
	}
}

func init() {
	extra.RegisterFuzzyDecoders()
	InitLog()
}
