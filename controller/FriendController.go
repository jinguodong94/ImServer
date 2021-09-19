package controller

import (
	"github.com/gin-gonic/gin"
	"imserver/dao"
	"imserver/req"
	"imserver/response"
)

type FriendController struct {
	BaseController
}

//添加好友
func (FriendController) AddFriend(ctx *gin.Context) {
	uid := GetUserIdFromToken(ctx)

	friendReq := &req.AddFriendReq{}
	ctx.ShouldBind(friendReq)

	relation := &dao.FriendRelation{}
	tx := dao.Db.Model(relation).Where("uid = ? and friend_id = ", uid, friendReq.FriendId).First(relation)
	if tx.Error == nil {
		if relation.Status == 1 {
			responseOk(ctx, response.NewErrorResponse("你已被对方拉黑，无法添加"))
			return
		} else if relation.Status == 0 {
			responseOk(ctx, response.NewErrorResponse("对方已是你好友"))
			return
		}
	}
	friendApply := &dao.FriendApply{}
	tx = dao.Db.Model(friendApply).Where("apply_uid = ? and to_uid = ?", uid, friendReq.FriendId).First(friendApply)
	if tx.Error == nil {
		tx = dao.Db.Model(friendApply).Updates(map[string]interface{}{"status": 0, "remarks": friendReq.Remarks})
	} else {
		friendApply.ApplyUid = uid
		friendApply.ToUid = friendReq.FriendId
		friendApply.Remarks = friendReq.Remarks
		tx = dao.Db.Create(friendApply)
	}
	if tx.Error != nil {
		responseOk(ctx, response.NewErrorResponse("好友申请发送失败"))
		return
	}
	responseOk(ctx, response.NewSuccessResponse("好友申请发送成功", nil))
}

//删除好友
func (FriendController) DelFriend(ctx *gin.Context) {
	uid := GetUserIdFromToken(ctx)
	deleteFriendReq := &req.DeleteFriendReq{}
	tx := dao.Db.Where(
		"uid = ? and friend_id = ?", uid, deleteFriendReq.OptionUid).Unscoped().Delete(&dao.FriendRelation{})
	if tx.Error != nil {
		responseOk(ctx, response.NewErrorResponse("删除失败"))
		return
	}
	responseOk(ctx, response.NewSuccessResponse("删除成功", deleteFriendReq.OptionUid))
}

//拉黑好友
func (FriendController) BlackFriend(ctx *gin.Context) {

}

//好友列表
func (FriendController) GetFriendList(ctx *gin.Context) {
	uid := GetUserIdFromToken(ctx)
	users := make([]dao.Users, 0, 10)
	dao.Db.Table(dao.TableName_FriendRelations+" as fr").Select(
		"u.*").Joins(
		"join users as u on u.id = fr.uid").Where(
		"uid = ? and status = 0", uid).Order("fr.updated_at").Find(&users)
	responseOk(ctx, response.NewSuccessResponse("查询成功", users))
}

//操作好友申请
func (FriendController) OptionFriendApply(ctx *gin.Context) {
	uid := GetUserIdFromToken(ctx)
	optionReq := &req.OptionFriendApplyReq{}
	ctx.ShouldBind(optionReq)
	tx := dao.Db.Model(&dao.FriendApply{}).Where(
		"to_uid = ? and apply_uid = ?", uid, optionReq.OptionId).Update("status", optionReq.OptionType)
	if tx.Error != nil {
		responseOk(ctx, response.NewErrorResponse("操作失败"))
		return
	}
	responseOk(ctx, response.NewSuccessResponse("操作成功", nil))
}
