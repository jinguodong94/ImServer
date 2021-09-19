package socket

import (
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"imserver/conf"
	"imserver/constant"
	"imserver/dao"
	"imserver/model"
	"imserver/mq"
	"imserver/utils"
	"log"
	"strconv"
	"time"
)

//处理心跳
func processHeartbeat(client *Client, heartbeat *model.Heartbeat) {
	log.Println("处理心跳", heartbeat.Token)
	uid, err := utils.TokenUtils.GetUserId(heartbeat.Token)
	if client.UserId == "" || err != nil {
		data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_heartbeat_error, Data: "心跳发送失败"})
		client.SendMessage(data)
		return
	}
	info := &model.OnLineUserInfo{}
	info, err = utils.RedisUtils.GetUserInfoById(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_heartbeat_error, Data: "未登录请先登录"})
		client.SendMessage(data)
		return
	}
	if info.Token != heartbeat.Token {
		data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_heartbeat_error, Data: "未登录请先登录"})
		client.SendMessage(data)
		return
	}
	info.HeartbeatExpires = utils.TimeUtils.GetCurrentTime()
	userInfoJson, _ := jsoniter.MarshalToString(info)
	dao.Rdb.HSet(context.Background(), constant.Redis_key_user_im_login_info, strconv.FormatUint(uint64(uid), 10), userInfoJson)
	data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_heartbeat_success, Data: "心跳发送成功"})
	client.SendMessage(data)
}

//处理登录
func processLogin(client *Client, login *model.Login) {
	log.Println("处理登录请求", login.Account)
	user := &dao.Users{}
	result := dao.Db.Model(user).Where("account = ? and pwd = ? and deleted_at is null", login.Account, login.Pwd).First(user)

	if result.Error != nil {
		data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_login_error, Data: "登录失败，账户或者密码错误"})
		client.SendMessage(data)
		return
	}
	//生成token
	token := utils.TokenUtils.CreateToken(user.ID)

	onlineUserInfo := &model.OnLineUserInfo{
		Token:            token,
		UserId:           user.ID,
		HeartbeatExpires: utils.TimeUtils.GetCurrentTime(),
		ServerId:         conf.Configs.ServerConfig.ServerId,
	}

	userInfoJson, _ := jsoniter.MarshalToString(onlineUserInfo)

	userId := strconv.FormatUint(uint64(user.ID), 10)
	rdbResult := dao.Rdb.HSet(context.Background(), constant.Redis_key_user_im_login_info, userId, userInfoJson)
	if rdbResult.Err() != nil {
		data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_login_error, Data: "登录异常"})
		client.SendMessage(data)
		return
	}
	client.UserId = userId
	ClientMgr.AddLoginClient(client)
	data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_login_success, Data: "登录成功"})
	client.SendMessage(data)

	//登录成功后检查是否有离线消息
	checkOffLineMessage(client.UserId)
}

//检查离线消息
func checkOffLineMessage(userId string) {
	messages := make([]dao.Messages, 0, 6)
	//单聊消息
	tx := dao.Db.Model(&dao.Messages{}).Where("to_uid = ? and deleted_at is null", userId).Order("created_at").Find(&messages)
	if tx.Error == nil {
		for _, message := range messages {
			mqMessage := &model.MqMessage{Data: &message}
			jsonMsg, _ := jsoniter.Marshal(mqMessage)
			processMqMessage(string(jsonMsg))
		}
	}
	//群聊消息
	groupMessages := make([]model.GroupMessage, 0, 6)
	tx = dao.Db.Table(dao.TableName_GroupMessagesRelation+" as gmr").Select("gms.*,gmr.uid as to_uid").Joins(
		fmt.Sprintf("join %s as gms on gmr.msg_id = gms.id", dao.TableName_GroupMessages)).Where(
		"gmr.uid = ?", userId).Find(&groupMessages)
	if tx.Error != nil {
		log.Println(tx.Error)
		return
	}
	for _, value := range groupMessages {
		mqMessage := &model.MqMessage{Type: 1, Data: &value}
		jsonMsg, _ := jsoniter.Marshal(mqMessage)
		processMqMessage(string(jsonMsg))
	}
}

//处理消息
func processMessage(client *Client, message *dao.Messages) {
	log.Println(fmt.Sprintf("处理发送消息 fromId:%s  toId:%s", message.FromUid, message.ToUid))

	if message.ChatType == 0 {
		//单聊
		processSingleChat(client, message)
	} else if message.ChatType == 1 {
		//群聊
		processGroupChat(client, message)
	} else if message.ChatType == 2 {
		//聊天室
		processRoomChat(client, message)
	}
}

//处理单聊
func processSingleChat(client *Client, message *dao.Messages) {
	//聊天室或者消息指令不存入表
	if message.ChatType != 2 && message.MsgType != 1 && message.IsOffLine == 0 {
		//存入数据库
		dao.Db.AutoMigrate(message)
		message.Status = 0
		tx := dao.Db.Create(message)
		if tx.Error != nil {
			data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_send_message_error, Data: message.MsgFlag})
			client.SendMessage(data)
			return
		}
	}
	data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_send_message_success, Data: message.MsgFlag})
	client.SendMessage(data)
	//获取发送者信息
	userInfo, err := utils.RedisUtils.GetUserInfoById(strconv.Itoa(int(message.ToUid)))
	if err != nil {
		//查询不到信息，说明对方没在线，就不管了，上线之后会发滞留的消息
		return
	}
	mqMessage := &model.MqMessage{Data: message}
	jsonMsg, _ := jsoniter.Marshal(mqMessage)
	if conf.Configs.ServerConfig.ServerId == userInfo.ServerId {
		processMqMessage(string(jsonMsg))
	} else {
		mq.MQ.SendMessage(jsonMsg, userInfo.ServerId)
	}
}

//处理群聊
func processGroupChat(client *Client, message *dao.Messages) {

	groupMessage := &dao.GroupMessages{}
	groupMessage.FromUid = message.FromUid
	groupMessage.ToGroupId = message.ToUid
	groupMessage.ExtendInfo = message.ExtendInfo
	groupMessage.MessageContent = message.MessageContent
	groupMessage.MsgType = message.MsgType
	groupMessage.IsOffLine = message.IsOffLine
	groupMessage.MessageType = message.MessageType

	if message.MsgType != 1 && message.IsOffLine == 0 {
		dao.Db.AutoMigrate(groupMessage)
		t := dao.Db.Create(groupMessage)
		if t.Error != nil {
			data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_send_message_error, Data: message.MsgFlag})
			client.SendMessage(data)
			return
		}
	}

	data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_send_message_success, Data: message.MsgFlag})
	client.SendMessage(data)

	//获取群内成员信息
	users := make([]dao.UserGroupRelation, 0, 50)
	tx := dao.Db.Model(&dao.UserGroupRelation{}).Where("group_id = ? and deleted_at is null", message.ToUid).Find(&users)
	if tx.Error != nil {
		//没查询到群内的人
		return
	}
	dao.Db.AutoMigrate(&dao.GroupMessagesRelation{})

	if message.MsgType != 1 && message.IsOffLine == 0 {
		start := time.Now() // 获取当前时间
		relations := make([]*dao.GroupMessagesRelation, 0, 50)
		//存入表
		for _, value := range users {
			if groupMessage.ID != 0 {
				relation := &dao.GroupMessagesRelation{}
				relation.Uid = value.Uid
				relation.MsgId = groupMessage.ID
				relations = append(relations, relation)
			}
		}
		dao.Db.Create(relations)
		elapsed := time.Since(start)
		log.Println("群消息sql插入执行时间", elapsed)
	}

	for _, value := range users {
		msg := &model.GroupMessage{
			ToUid:         value.Uid,
			GroupMessages: *groupMessage,
		}
		mqMessage := &model.MqMessage{Type: 1, Data: msg}
		jsonMsg, _ := jsoniter.Marshal(mqMessage)

		cli := ClientMgr.GetClientByUserId(strconv.Itoa(int(value.Uid)))
		if cli != nil {
			processMqMessage(string(jsonMsg))
		} else {
			userInfo, err := utils.RedisUtils.GetUserInfoById(strconv.Itoa(int(value.Uid)))
			if err != nil {
				continue
			}
			mq.MQ.SendMessage(jsonMsg, userInfo.ServerId)
		}
	}
}

//处理聊天室
func processRoomChat(client *Client, message *dao.Messages) {
	members, err := utils.RedisUtils.GetRoomMembers(message.ToUid)
	if err != nil {
		data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_send_message_error, Data: message.MsgFlag})
		client.SendMessage(data)
		log.Println("聊天室不存在 > ", message.ToUid)
		return
	}
	for _, uid := range members {
		u, _ := strconv.Atoi(uid)
		message.ToUid = uint(u)
		mqMessage := &model.MqMessage{Type: 2, Data: message}
		jsonMsg, _ := jsoniter.Marshal(mqMessage)

		cli := ClientMgr.GetClientByUserId(uid)
		if cli != nil {
			processMqMessage(string(jsonMsg))
		} else {
			userInfo, err := utils.RedisUtils.GetUserInfoById(uid)
			if err != nil {
				continue
			}
			mq.MQ.SendMessage(jsonMsg, userInfo.ServerId)
		}
	}
	data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_send_message_success, Data: message.MsgFlag})
	client.SendMessage(data)
}

//处理消息确认
func processAckMessage(client *Client, ackMessage *model.AckMessage) {
	defer func() {
		//消息回应
		data, _ := jsoniter.Marshal(model.MsgResponse{RespType: constant.Resp_ack_message_success, Data: ackMessage.MessageId})
		client.SendMessage(data)
	}()

	if ackMessage.ChatType == 0 {
		messages := &dao.Messages{}
		tx := dao.Db.Model(messages).Where("id = %s", ackMessage.MessageId).First(messages)
		if tx.Error != nil {
			//没查询到消息不管了
			return
		}
		begin := dao.Db.Begin()
		//删除临时消息表
		tx = begin.Unscoped().Delete(messages)
		if tx.Error != nil {
			begin.Rollback()
			return
		}
		//新增到消息记录表
		messageHistory := &dao.MessageHistory{}
		begin.AutoMigrate(messageHistory)
		messageHistory.ChatType = messages.ChatType
		messageHistory.FromUid = messages.FromUid
		messageHistory.ToUid = messages.ToUid
		messageHistory.ExtendInfo = messages.ExtendInfo
		messageHistory.MessageContent = messages.MessageContent
		tx = begin.Create(messageHistory)
		if tx.Error != nil {
			begin.Rollback()
			return
		}
		begin.Commit()
	} else if ackMessage.ChatType == 1 {
		message := &dao.GroupMessagesRelation{}
		tx := dao.Db.Model(message).Where("msg_id = %s and uid = ?", ackMessage.MessageId, client.UserId).First(message)
		if tx.Error != nil {
			return
		}
		begin := dao.Db.Begin()
		//删除临时消息表
		tx = begin.Unscoped().Delete(message)
		if tx.Error != nil {
			begin.Rollback()
			return
		}
		messageHistory := &dao.GroupMessagesHistoryRelation{}
		messageHistory.Uid = message.Uid
		messageHistory.MsgId = message.MsgId
		tx = begin.Create(messageHistory)
		if tx.Error != nil {
			begin.Rollback()
			return
		}
		begin.Commit()
	}
}
