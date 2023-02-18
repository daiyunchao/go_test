package server_exit

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

// 服务器优雅退出
func TestServerExit(t *testing.T) {
	app := NewApp()
	server1 := HttpServer{
		Point: "8001",
		Handle: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello 8001!\n")
		},
		StartedHandle: func() error {
			fmt.Println("Hello 8001!")
			fmt.Println("在这里假装连接数据库等初始化工作")
			return nil
		},
		StopHandle: func() error {
			fmt.Println("Bye 8001!")
			fmt.Println("断开数据库连接,释放资源")
			return nil
		},
	}
	app.RegisterServer(&server1)

	server2 := TcpServer{
		Point: "8002",
		Handle: func(conn net.Conn) {
			defer conn.Close()
			fmt.Println("处理请求")
		},
		StartedHandle: func() error {
			fmt.Println("Hello 8002!")
			fmt.Println("在这里假装连接数据库等初始化工作")
			return nil
		},
		StopHandle: func() error {
			fmt.Println("Bye 8002!")
			fmt.Println("断开数据库连接,释放资源")
			return nil
		},
	}
	app.RegisterServer(&server2)

	app.AppStart()
}

// App
type App struct {
	signalCount int
	signalChan  chan os.Signal //信号
	startWg     sync.WaitGroup
	stopWg      sync.WaitGroup
	servers     []IServer
}

func NewApp() *App {
	return &App{
		signalCount: 0,
		signalChan:  make(chan os.Signal),
		servers:     make([]IServer, 0),
	}
}
func (app *App) AppStart() {
	if len(app.servers) == 0 {
		panic("no server in app")
	}
	fmt.Println("服务器启动前,可以开始做一些公共初始化的事情")
	for _, server := range app.servers {
		app.startWg.Add(1)
		go server.Start(&app.startWg)
	}
	app.startWg.Wait()
	fmt.Println("全部服务器都已完成启动,可以开始做一些服务器启动后的事情")

	//开始监听退出信号
	signal.Notify(app.signalChan, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		select {
		case <-app.signalChan:
			app.signalCount += 1
			if app.signalCount > 1 {
				//强行退出
				return
			}
			//如果是第一次退出,则释放资源
			//额外启动一个timer,如果时间到了还未退出则强行退出
			time.AfterFunc(time.Second*30, func() {
				fmt.Println("超时了,退出")
				return
			})
			for _, server := range app.servers {
				app.stopWg.Add(1)
				go server.Stop(&app.stopWg)
			}
			app.stopWg.Wait()
			fmt.Println("资源释放完毕可以退出了")
			return
		}
	}
}

func (app *App) RegisterServer(s IServer) {
	app.servers = append(app.servers, s)
}

// 服务器接口
type IServer interface {
	Start(wg *sync.WaitGroup)
	Stop(wg *sync.WaitGroup)
}

type TcpServer struct {
	Point             string //启动端口
	Handle            func(conn net.Conn)
	StartedHandle     func() error
	StopHandle        func() error
	canProcessRequest bool
}

func (s *TcpServer) Start(wg *sync.WaitGroup) {
	s.canProcessRequest = true
	fmt.Println("开始启动", s.Point, "服务器")
	listen, err := net.Listen("tcp", fmt.Sprint("0.0.0.0:", s.Point))
	if err != nil {
		fmt.Println("启动服务器", s.Point, "出现异常", err)
		panic("启动服务器异常")
		return
	}
	err = s.StartedHandle()
	if err != nil {
		panic("启动服务器异常")
		return
	}
	wg.Done()
	fmt.Println("完成启动", s.Point, "服务器")

	for {
		//2.接收客户端的链接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			continue
		}
		//3.开启一个Goroutine，处理链接
		if s.canProcessRequest {
			go s.Handle(conn)
		}
	}
}

func (s *TcpServer) Stop(wg *sync.WaitGroup) {
	fmt.Println("开始停止", s.Point, "服务器")
	s.canProcessRequest = false
	fmt.Println("5秒后", s.Point, "服务器将会被关闭")
	time.Sleep(time.Second * 5)
	//5秒后关闭服务
	err := s.StopHandle()
	if err != nil {
		fmt.Println("停止服务器", s.Point, "出现异常", err)
		panic("停止服务器异常")
		return
	}
	wg.Done()
	fmt.Println("完成停止", s.Point, "服务器")
}

// 服务器实现
type HttpServer struct {
	Point             string //启动端口
	Handle            func(writer http.ResponseWriter, request *http.Request)
	StartedHandle     func() error
	StopHandle        func() error
	canProcessRequest bool
}

func (s *HttpServer) Start(wg *sync.WaitGroup) {
	fmt.Println("开始启动", s.Point, "服务器")
	s.canProcessRequest = true
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if s.canProcessRequest {
			s.Handle(w, r)
		} else {
			//不能再执行新的请求了,返回503
		}
	})
	go func() {
		err := http.ListenAndServe(fmt.Sprint(":", s.Point), nil)
		if err != nil {
			fmt.Println("启动服务器", s.Point, "出现异常", err)
			panic("启动服务器异常")
			return
		}
	}()
	time.Sleep(1 * time.Second)
	err := s.StartedHandle()
	if err != nil {
		panic("启动服务器异常")
		return
	}
	wg.Done()
	fmt.Println("完成启动", s.Point, "服务器")
}

func (s *HttpServer) Stop(wg *sync.WaitGroup) {
	fmt.Println("开始停止", s.Point, "服务器")
	s.canProcessRequest = false
	fmt.Println("5秒后", s.Point, "服务器将会被关闭")
	time.Sleep(time.Second * 5)
	err := s.StopHandle()
	if err != nil {
		fmt.Println("停止服务器", s.Point, "出现异常", err)
		panic("停止服务器异常")
		return
	}
	wg.Done()
	fmt.Println("完成停止", s.Point, "服务器")
}
