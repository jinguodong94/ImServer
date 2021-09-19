package socket

import (
	"sync"
)

var (
	ClientMgr    = NewClientManager()
	clientRwLock sync.RWMutex
)

type ClientManager struct {
	//所有客户端
	allClient map[*Client]bool
	//已登录的客户端
	loginedClient map[string]*Client

	addClientChan      chan *Client
	addLoginClientChan chan *Client
	removeClientChan   chan *Client
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		allClient:          make(map[*Client]bool),
		loginedClient:      make(map[string]*Client),
		addClientChan:      make(chan *Client, 10),
		addLoginClientChan: make(chan *Client, 10),
		removeClientChan:   make(chan *Client, 10),
	}
	return
}

func (mgr *ClientManager) Init() {
	go func() {
		for {
			select {
			case client := <-ClientMgr.addClientChan:
				mgr.allClient[client] = true

			case client := <-ClientMgr.addLoginClientChan:
				mgr.loginedClient[client.UserId] = client

			case client := <-ClientMgr.removeClientChan:
				delete(mgr.allClient, client)
				if client.UserId != "" {
					delete(mgr.loginedClient, client.UserId)
				}
			}
		}
	}()
}

//添加客户端到所有
func (mgr *ClientManager) AddClient(client *Client) {
	mgr.addClientChan <- client
}

//删除客户端
func (mgr *ClientManager) DelClient(client *Client) {
	mgr.removeClientChan <- client
}

//添加已登录的客户端
func (mgr *ClientManager) AddLoginClient(client *Client) {
	mgr.addLoginClientChan <- client
}

//获取客户端
func (mgr *ClientManager) GetClientByUserId(userId string) *Client {
	clientRwLock.RLock()
	defer clientRwLock.RUnlock()
	return mgr.loginedClient[userId]
}

//获取所有客户端
func (mgr *ClientManager) GetAllClient() map[*Client]bool {
	return mgr.allClient
}

//获取已登录的客户端
func (mgr *ClientManager) GetLoginClient() map[string]*Client {
	return mgr.loginedClient
}
