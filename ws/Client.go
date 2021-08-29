package ws

import (
	"errors"
	"fmt"
	"gindemo/constant"
	"gindemo/model"
	"gindemo/utils"
	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
	"log"
)

type Client struct {
	IpAddress     string          //IP地址
	UserId        string          //用户id
	HeartbeatTime uint64          //心跳时间
	Socket        *websocket.Conn //客户端连接
	WriteData     chan []byte     //待写入数据
}

func NewClient(IpAddress string, Socket *websocket.Conn) (client *Client) {
	client = &Client{
		IpAddress:     IpAddress,
		HeartbeatTime: utils.TimeUtils.GetCurrentTime(),
		Socket:        Socket,
		WriteData:     make(chan []byte, 50),
	}
	return
}

//开启客户端
func (client *Client) Open() {
	log.Println("client open ip is :", client.IpAddress)
	go client.loopRead()
	go client.loopWrite()
	ClientMgr.AddClient(client)
}

//循环读取
func (client *Client) loopRead() {

	defer func() {
		log.Println("close client channel", GetClientString(client))
		close(client.WriteData)
	}()

	for {
		_, message, err := client.Socket.ReadMessage()

		if err != nil {
			log.Println("client read error", GetClientString(client), "errInfo->", err)
			return
		}

		client.processData(message)
	}
}

//循环写入
func (client *Client) loopWrite() {

	clientString := GetClientString(client)

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("process data error :", rec)
		}
	}()

	defer func() {
		log.Println("client close connect", clientString)
		//删除客户端并关闭连接
		ClientMgr.DelClient(client)
		client.Socket.Close()
	}()

	for {
		data, ok := <-client.WriteData
		if !ok {
			log.Println("client write channel close", clientString)
			return
		}

		err := client.Socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("client write error", clientString)
			break
		}
		log.Println(clientString, "send message success : ", string(data))
	}
}

//发送消息
func (client *Client) sendMessage(data []byte) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.New("channel is close")
		}
	}()
	client.WriteData <- data
	return
}

//处理数据
func (client *Client) processData(message []byte) {

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("process data error :", rec)
		}
	}()

	messageData := string(message)
	log.Println("process data", GetClientString(client), " data -> ", messageData)

	request, err := GetRequest(messageData, nil)
	if err != nil {
		log.Println("json deserialization err -> ", err)
	}
	switch constant.ReqType(request.Type) {

	case constant.Req_heartbeat:
		//心跳
		request, err := GetRequest(messageData, &model.Heartbeat{})
		if err != nil {
			log.Println("json deserialization err -> ", err)
			break
		}
		heartbeat := request.Data.(*model.Heartbeat)
		processHeartbeat(client, heartbeat)
	case constant.Req_login:
		//登录
		request, err := GetRequest(messageData, &model.Login{})
		if err != nil {
			log.Println("json deserialization err -> ", err)
			break
		}
		login := request.Data.(*model.Login)
		processLogin(client, login)
	case constant.Req_sendMessage:
		//发送消息
		request, err := GetRequest(messageData, &model.Message{})
		if err != nil {
			log.Println("json deserialization err -> ", err)
			break
		}
		message := request.Data.(*model.Message)
		processMessage(client, message)
	}
}

func GetRequest(jsonData string, data interface{}) (request *model.MsgRequest, err error) {
	request = &model.MsgRequest{Data: data}
	err = jsoniter.UnmarshalFromString(jsonData, request)
	return
}

func GetClientString(client *Client) (str string) {
	str = fmt.Sprintf("[ip : %s] [userId : %s] ", client.IpAddress, client.UserId)
	return
}
