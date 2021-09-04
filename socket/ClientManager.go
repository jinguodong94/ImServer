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
	allClient []*Client
	//已登录的客户端
	loginedClient map[string]*Client
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		allClient:     make([]*Client, 0, 10),
		loginedClient: make(map[string]*Client),
	}
	return
}

//添加客户端到所有
func (mgr *ClientManager) AddClient(client *Client) {
	clientRwLock.Lock()
	defer clientRwLock.Unlock()
	mgr.allClient = append(mgr.allClient, client)
}

//删除客户端
func (mgr *ClientManager) DelClient(client *Client) {
	clientRwLock.Lock()
	defer clientRwLock.Unlock()
	remove(mgr.allClient, client)
	mgr.allClient = remove(mgr.allClient, client)
	if client.UserId != "" {
		delete(mgr.loginedClient, client.UserId)
	}
}

//添加已登录的客户端
func (mgr *ClientManager) AddLoginClient(client *Client) {
	clientRwLock.Lock()
	defer clientRwLock.Unlock()
	mgr.loginedClient[client.UserId] = client
}

//获取客户端
func (mgr *ClientManager) GetClientByUserId(userId string) *Client {
	clientRwLock.RLock()
	defer clientRwLock.RUnlock()
	return mgr.loginedClient[userId]
}

//获取所有客户端
func (mgr *ClientManager) GetAllClient() []*Client {
	return mgr.allClient
}

//获取已登录的客户端
func (mgr *ClientManager) GetLoginClient() map[string]*Client {
	return mgr.loginedClient
}

func remove(slice []*Client, e *Client) []*Client {
	position := getPosition(slice, e)
	return removeByPosition(slice, position)
}

func removeByPosition(slice []*Client, position int) []*Client {
	if position != -1 {
		return append(slice[:position], slice[position+1:]...)
	}
	return slice
}

//获取元素位置
func getPosition(slice []*Client, e *Client) int {
	for index, value := range slice {
		if value == e {
			return index
		}
	}
	return -1
}
