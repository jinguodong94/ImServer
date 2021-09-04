package socket

import (
	"gindemo/constant"
	"gindemo/dao"
	"gindemo/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

//mq消息接收
func Consumer(d *amqp.Delivery) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Consumer process data error :", rec)
		}
	}()
	defer d.Ack(true)
	jsonString := string(d.Body)
	log.Println("mq接收消息:", jsonString, "   routeKey", d.RoutingKey)
	processMqMessage(jsonString)
}

func processMqMessage(jsonString string) {
	mqMessage := &model.MqMessage{}
	err := jsoniter.UnmarshalFromString(jsonString, mqMessage)
	if err != nil {
		log.Println("json解析失败 : ", jsonString)
		return
	}

	if mqMessage.Type == 0 {
		processChatMessage(jsonString)
	} else if mqMessage.Type == 1 {
		processGroupMessage(jsonString)
	}
}

func processGroupMessage(jsonString string) {
	log.Println("处理群聊 : ", jsonString)
	mqMessage := &model.MqMessage{Data: &model.GroupMessage{}}
	err := jsoniter.UnmarshalFromString(jsonString, mqMessage)
	if err != nil {
		log.Println("json解析失败 : ", jsonString)
		return
	}
	message := mqMessage.Data.(*model.GroupMessage)
	toClient := ClientMgr.GetClientByUserId(strconv.Itoa(int(message.ToUid)))
	if toClient != nil {
		msg := &dao.Messages{}
		msg.ChatType = 1
		msg.MessageType = message.MessageType
		msg.FromUid = message.FromUid
		msg.ToUid = message.ToGroupId
		msg.ID = message.ID
		msg.ExtendInfo = message.ExtendInfo
		msg.MessageContent = message.MessageContent
		messageBytes, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_recive_message, Data: msg})
		err := toClient.SendMessage(messageBytes)
		if err != nil {
			return
		}
		if message.MsgType != 1 && message.IsOffLine == 0 {
			dao.Db.Model(&dao.GroupMessagesRelation{}).Where("msg_id = ? and uid = ?", message.ID, message.ToUid).Update("status", 1)
		}
	}
}

//处理聊天消息
func processChatMessage(jsonString string) {
	log.Println("处理单聊 : ", jsonString)
	mqMessage := &model.MqMessage{Data: &dao.Messages{}}
	err := jsoniter.UnmarshalFromString(jsonString, mqMessage)
	if err != nil {
		log.Println("json解析失败 : ", jsonString)
		return
	}
	message := mqMessage.Data.(*dao.Messages)
	toClient := ClientMgr.GetClientByUserId(strconv.Itoa(int(message.ToUid)))
	if toClient != nil {
		messageBytes, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_recive_message, Data: message})
		err = toClient.SendMessage(messageBytes)
		if err != nil {
			return
		}
		if message.ChatType != 2 && message.MsgType != 1 && message.IsOffLine == 0 {
			dao.Db.Model(message).Update("status", 1)
		}
	}
	return
}
