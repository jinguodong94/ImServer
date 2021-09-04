package route

import (
	"gindemo/controller"
	"github.com/gin-gonic/gin"
	"log"
)

var Route *gin.Engine

func Init() {

	log.Println("初始化路由")

	Route = gin.Default()

	//用户相关
	userGroup := Route.Group("/user")
	{
		//登录
		userGroup.POST("/login", controller.UserController{}.Login)
		//注册
		userGroup.POST("/register", controller.UserController{}.Register)
		//获取个人信息
		userGroup.POST("/getUserInfo", controller.UserController{}.GetUserInfo)
		//修改个人信息
		userGroup.POST("/updateUserInfo", controller.UserController{}.UpdateUserInfo)
	}

	//好友相关
	friendGroup := Route.Group("/friend")
	{
		//添加好友
		friendGroup.POST("/addFriend", controller.FriendController{}.AddFriend)
		//删除好友
		friendGroup.POST("/delFriend", controller.FriendController{}.DelFriend)
		//拉黑好友
		friendGroup.POST("/blackFriend", controller.FriendController{}.BlackFriend)
		//好友列表
		friendGroup.POST("/getFriendList", controller.FriendController{}.GetFriendList)
	}

	//群相关
	groupGroup := Route.Group("/group")
	groupGroup.Use(controller.AuthController{}.LoginAuth)
	{
		groupGroup.POST("/createGroup", controller.GroupController{}.CreateGroup)

		groupGroup.POST("/joinGroup", controller.GroupController{}.JoinGroup)

		groupGroup.POST("/exitGroup", controller.GroupController{}.ExitGroup)

		groupGroup.POST("/dissolveGroup", controller.GroupController{}.DissolveGroup)

		groupGroup.GET("/getGroupList", controller.GroupController{}.GetGroupList)
	}

	Route.GET("/systemStatus", controller.SystemController{}.GetSystemStatus)
}
