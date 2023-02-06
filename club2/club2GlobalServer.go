package club2

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type GlobalServer struct {
	sync.Mutex   //资源锁
	reqChan      chan Request
	reqChanClose chan bool
	resChan      chan Response
	timeoutMap   map[int]bool //记录超时
}

// CreateGlobalClubServer 创建一个全局的工会协程
func CreateGlobalClubServer() *GlobalServer {
	server := &GlobalServer{
		reqChan:      make(chan Request),
		reqChanClose: make(chan bool),
		resChan:      make(chan Response),
	}
	go server.Run()
	return server
}

func (gs *GlobalServer) Run() {
	for {
		select {
		case req := <-gs.reqChan:
			gs.Lock()
			//处理正常全局的业务逻辑
			log.Print(req)
			gs.Unlock()
		}
	}
}

func (gs *GlobalServer) SendRequest(req Request) (Response, error) {
	//简单验证Request
	gs.reqChan <- req
	var isResponse = false
	for {
		select {
		case res := <-gs.resChan:
			isResponse = true
			return res, nil
		case <-time.After(time.Duration(time.Second * 3)):
			if !isResponse {
				fmt.Println("timeout")
				//标识该请求已经超时了
				gs.timeoutMap[req.reqId] = true
				return Response{}, errors.New("timeout")
			}
		}
	}
}
