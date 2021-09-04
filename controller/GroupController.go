package controller

import (
	"fmt"
	"gindemo/dao"
	"gindemo/req"
	"gindemo/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type GroupController struct {
	BaseController
}

//创建群
func (GroupController) CreateGroup(ctx *gin.Context) {
	createReq := &req.CreateGroupReq{}
	err := ctx.ShouldBind(createReq)
	if err != nil {
		responseOk(ctx, response.NewErrorResponse("参数有误"))
		return
	}
	uid := GetUserIdFromToken(ctx)

	if createReq.GroupName == "" {
		createReq.GroupName = fmt.Sprintf("%s_%s", "group", strconv.Itoa(int(uid)))
	}

	groups := &dao.Groups{}
	dao.Db.AutoMigrate(groups)

	//开启事务
	dao.Db.Begin()

	groups.GroupName = createReq.GroupName
	groups.Icon = createReq.Icon
	tx := dao.Db.Create(groups)
	if tx.Error != nil {
		responseOk(ctx, response.NewErrorResponse("创建失败"))
		dao.Db.Rollback()
		return
	}
	//添加房主到群关系表
	userGroupRelation := &dao.UserGroupRelation{}
	dao.Db.AutoMigrate(userGroupRelation)
	userGroupRelation.Uid = uid
	userGroupRelation.GroupId = groups.ID
	userGroupRelation.Role = 1
	tx = dao.Db.Create(userGroupRelation)
	if tx.Error != nil {
		responseOk(ctx, response.NewErrorResponse("创建失败"))
		dao.Db.Rollback()
		return
	}
	dao.Db.Commit()
	responseOk(ctx, response.NewSuccessResponse("创建成功", groups))
}

//加群
func (GroupController) JoinGroup(ctx *gin.Context) {
	groupReq := &req.JoinOrExitGroupReq{}
	err := ctx.ShouldBind(groupReq)
	if err != nil {
		responseOk(ctx, response.NewErrorResponse("参数有误"))
		return
	}

	uid := GetUserIdFromToken(ctx)
	userGroupRelation := &dao.UserGroupRelation{}
	dao.Db.AutoMigrate(userGroupRelation)
	tx := dao.Db.Model(userGroupRelation).Where("group_id = ? and uid = ?", groupReq.GroupId, uid).First(userGroupRelation)
	if tx.Error == nil {
		responseOk(ctx, response.NewErrorResponse("不可重复加入此群"))
		return
	}

	userGroupRelation.Uid = uid
	groupId, _ := strconv.Atoi(groupReq.GroupId)
	userGroupRelation.GroupId = uint(groupId)
	userGroupRelation.Role = 0
	tx = dao.Db.Create(userGroupRelation)
	if tx.Error != nil {
		responseOk(ctx, response.NewErrorResponse("加入失败"))
		return
	}
	responseOk(ctx, response.NewSuccessResponse("加入成功", userGroupRelation))
}

//退群
func (GroupController) ExitGroup(ctx *gin.Context) {
	groupReq := &req.JoinOrExitGroupReq{}
	err := ctx.ShouldBind(groupReq)
	if err != nil {
		responseOk(ctx, response.NewErrorResponse("参数有误"))
		return
	}

	uid := GetUserIdFromToken(ctx)
	userGroupRelation := &dao.UserGroupRelation{}
	dao.Db.AutoMigrate(userGroupRelation)
	tx := dao.Db.Where("group_id = ? and uid = ?", groupReq.GroupId, uid).Unscoped().Delete(userGroupRelation)
	if tx.Error != nil {
		responseOk(ctx, response.NewErrorResponse("退群失败"))
		return
	}
	responseOk(ctx, response.NewSuccessResponse("退群成功", nil))
}

//解散群
func (GroupController) DissolveGroup(ctx *gin.Context) {

}

//获取群列表
func (GroupController) GetGroupList(ctx *gin.Context) {
	uid := GetUserIdFromToken(ctx)
	groups := make([]dao.Groups, 0, 1)
	dao.Db.Table(dao.TableName_UserGroupRelations+" as ugr").Select(
		"gp.*,ugr.created_at as joinTime").Joins(
		"join groups as gp on ugr.group_id = gp.id").Where("ugr.uid = ?", uid).Find(&groups)
	responseOk(ctx, response.NewSuccessResponse("查询成功", groups))
}
