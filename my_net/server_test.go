package my_net

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
)

const (
	MSG_TYPE_HANDSHAKE = 1
)

func TestNetServer(t *testing.T) {
	server := Server{
		port: "9001",
	}
	server.Start()
}

type Server struct {
	port string
}

// 启动服务器
func (s *Server) Start() {
	listen, err := net.Listen("tcp", "127.0.0.1:"+s.port)
	if err != nil {
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			return
		}

		//每一个连接一个协程去处理
		go s.handler(conn)
	}
}

// 规定包的格式
// |消息类型2字节|路由8字节|消息长度4字节|消息|
func (s *Server) handler(conn net.Conn) {
	for {
		//消息类型
		var msgTypeBuffer = make([]byte, 2)
		_, err := conn.Read(msgTypeBuffer)
		if err != nil {
			conn.Close()
			return
		}
		msgType := binary.BigEndian.Uint16(msgTypeBuffer)
		fmt.Println("rev msgType:", msgType)

		//路由
		var msgRouteBuffer = make([]byte, 8)
		_, err = conn.Read(msgRouteBuffer)
		if err != nil {
			conn.Close()
			return
		}
		msgRoute := binary.BigEndian.Uint64(msgRouteBuffer)
		fmt.Println("rev msgRoute:", msgRoute)

		//消息长度
		var msgLenBuffer = make([]byte, 4)
		_, err = conn.Read(msgLenBuffer)
		if err != nil {
			conn.Close()
			return
		}
		var msgLen = binary.BigEndian.Uint32(msgLenBuffer)
		fmt.Println("rev msgLen:", msgLen)

		//消息体
		var msgBody = make([]byte, msgLen)
		_, err = conn.Read(msgBody)
		if err != nil {
			conn.Close()
			return
		}
		fmt.Printf("rev msgBody: %s\n", msgBody)
		//返回数据
		conn.Write([]byte("hello"))
	}
}
