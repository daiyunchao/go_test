package club

import (
	"fmt"
	"sync"
)

const (
	EventUserLogin = 1
)

type Club struct {
	Id    int
	Name  string
	Score int
}
type Event struct {
	*Club
	EventType int
}

type Server struct {
	inClubChan    chan Event
	closeClubChan chan bool  //关闭Channel
	stopClubChan  chan bool  //停止Channel
	isClose       bool       //是否关闭
	isStop        bool       //是否停止
	closeMutex    sync.Mutex //关闭锁
	stopMutex     sync.Mutex //停止锁
}

func (server *Server) Init() {
	server.inClubChan = make(chan Event)
	server.closeClubChan = make(chan bool)
	server.stopClubChan = make(chan bool)
	server.isClose = false
	server.isStop = false
	server.closeMutex = sync.Mutex{}
	server.stopMutex = sync.Mutex{}
}
func (server *Server) Open() {
	for {
		select {
		case clubEvent := <-server.inClubChan:
			//如果关闭了,不在接收新的请求
			if !server.isClose {
				if clubEvent.EventType == EventUserLogin {
					addClubScore(clubEvent.Club)
				}
			} else {
				fmt.Println("服务已关闭,不能在接收新的请求")
			}
		case _, ok := <-server.closeClubChan:
			if !ok {
				fmt.Println("服务关闭中")
				//不再读取
				server.closeClubChan = nil
			}
		case _, ok := <-server.stopClubChan:
			if !ok {
				fmt.Println("服务停止中")
				//不再读取
				server.stopClubChan = nil
			}
			break
		}
	}
}

// CloseServer 关闭服务,服务不在接收新的请求
func (server *Server) CloseServer() {
	server.closeMutex.Lock()
	if !server.isClose {
		server.isClose = true
		close(server.closeClubChan)
	}
	server.closeMutex.Unlock()
}

// ReOpen 重新打开服务
func (server *Server) ReOpen() {
	if !server.isStop && server.isClose {
		//如果服务没有终止,直接打开
		server.isClose = false
		server.closeClubChan = make(chan bool)
	} else if server.isStop {
		//服务已终止,重新创建协程
		server.Init()
		go server.Open()
	} else {
		fmt.Println("不能ReOpen,状态异常")
	}
}

// StopServer 停止服务,协程退出
func (server *Server) StopServer() {
	server.stopMutex.Lock()
	if server.isClose && !server.isStop {
		server.isStop = true
		close(server.stopClubChan)
	}
	server.stopMutex.Unlock()
}

func (server *Server) InChan(clubEvent Event) {
	if !server.isStop {
		server.inClubChan <- clubEvent
	}
}
func addClubScore(club *Club) {
	club.Score += 1
	fmt.Printf("Club :%s current Score = %d\n", club.Name, club.Score)
}
