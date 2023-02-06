package club2

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Request 请求对象
type Request struct {
	reqId     int
	reqRoute  int
	reqAction map[string]interface{}
}

// Response 响应对象
type Response struct {
	reqId   int
	resId   int
	code    int
	resData map[string]interface{}
}

type Server struct {
	club         *Club2
	sync.Mutex   //资源锁
	reqChan      chan Request
	reqChanClose chan bool
	resChan      chan Response
	timeoutMap   map[int]bool //记录超时
}

type Club2 struct {
}

func (club *Club2) CreateClub(req Request) Response {
	res := Response{}
	log.Println("in CreateClub")
	time.Sleep(time.Second * 2)
	return res
}

// CreateClubServer 创建工会服务
// 当创建一个公会后立即创建一个工会对应的Server
func CreateClubServer(club *Club2) *Server {
	server := &Server{
		club:         club,
		reqChan:      make(chan Request),
		reqChanClose: make(chan bool),
		resChan:      make(chan Response),
	}
	go server.Run()
	return server
}

// Run 启动一个background服务
func (s *Server) Run() {
	for {
		select {
		case req := <-s.reqChan:
			s.Lock()
			//处理正常的业务逻辑
			res := s.club.CreateClub(req)
			if !s.timeoutMap[req.reqId] {
				s.resChan <- res
			} else {
				//如果超时了,不在放入chan中,避免chan堵塞
				log.Println("该请求已超时,不能再放入chan中")
				delete(s.timeoutMap, req.reqId)
			}
			s.Unlock()
		}
	}
}

func (s *Server) SendRequest(req Request) (Response, error) {
	//简单验证Request
	s.reqChan <- req
	var isResponse = false
	for {
		select {
		case res := <-s.resChan:
			isResponse = true
			return res, nil
		case <-time.After(time.Duration(time.Second * 3)):
			if !isResponse {
				fmt.Println("timeout")
				//标识该请求已经超时了
				s.timeoutMap[req.reqId] = true
				return Response{}, errors.New("timeout")
			}
		}
	}
}
