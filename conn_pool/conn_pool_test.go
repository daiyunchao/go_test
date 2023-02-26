package conn_pool

import (
	"net"
	"testing"
	"time"
)

// 没有连接池版本
func TestNoConnPool(t *testing.T) {
	server := Server{}
	server.Start()
	time.Sleep(time.Second * 2)
	//普通版本没有连接池,开启10个协程,创建10个链接
	for i := 0; i < 10; i++ {
		go func() {
			client := Client{}
			client.Start()
		}()
	}
}

// 有连接池的版本
func TestConnPool(t *testing.T) {

}

type ConnPool struct {
	maxConnCount int //最大链接数量
	connList     chan net.Conn
	acvCount     int //活跃数量
}

func CreateConnPool() *ConnPool {
	connPool := &ConnPool{
		maxConnCount: 5,
		connList:     make(chan net.Conn, 5),
		acvCount:     0,
	}

	return connPool
}

// 放回一个链接
func (c *ConnPool) Put(conn net.Conn) {
	if cap(c.connList) >= c.maxConnCount {
		//如果超过了数量,则直接丢弃
		conn.Close()
	}
	//如果有位置就放入
	c.acvCount -= 1
	c.connList <- conn
}

// 获得一个链接
func (c *ConnPool) Get(createConn func() net.Conn) net.Conn {
	if c.acvCount == 0 {
		//没有活跃的就创建一个
		newConn := createConn()
		c.connList <- newConn
	}
	c.acvCount += 1
	return <-c.connList
}

type Server struct {
}

// 启动服务器
func (s *Server) Start() {
	listen, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			return
		}
		go s.handler(conn)
	}
}
func (s *Server) handler(conn net.Conn) {
	//处理请求
	for {
		buffer := make([]byte, 8)
		conn.Read(buffer)
	}
}

type Client struct {
}

func (c *Client) Start() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		return
	}
	_, err = conn.Write([]byte("Hello"))
	if err != nil {
		conn.Close()
		return
	}
	for {
		buffer := make([]byte, 8)
		_, err := conn.Read(buffer)
		if err != nil {
			return
		}
	}
}
