package conf

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
)

type Config struct {
	MysqlConfig    MysqlConfig    `ini:"mysql"`
	RedisConfig    RedisConfig    `ini:"redis"`
	RabbitMQConfig RabbitMQConfig `ini:"rabbitMQ"`
	ServerConfig   ServerConfig   `ini:"server"`
}

var Configs = new(Config)

func Init() {

	log.Println("初始化配置")

	err := ini.MapTo(Configs, "conf/config.ini")
	if err != nil {
		panic(fmt.Sprint("init config file is error ", err.Error()))
		return
	}
}

type MysqlConfig struct {
	Address string `ini:"address"`
	Pwd     string `ini:"pwd"`
	User    string `ini:"user"`
	Db      string `ini:"db"`
}

type RedisConfig struct {
	Address string `ini:"address"`
	Pwd     string `ini:"pwd"`
}
type RabbitMQConfig struct {
	Address string `ini:"address"`
	User    string `ini:"user"`
	Pwd     string `ini:"pwd"`
}

type ServerConfig struct {
	ServerId string `ini:"serverId"`
}
