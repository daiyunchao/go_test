package main

import (
	"go_test/club"
	"go_test/club2"
	"log"
	"time"
)

// 协程+channel模拟: 同一时刻,100位玩家登录导致同一个工会的积分变化
func clubTest() {
	//启动一个协程,用于处理工会
	clubServer := club.Server{}
	clubServer.Init()
	go clubServer.Open()

	//模拟100个人登录
	clubInfo := club.Club{
		Id:    1,
		Name:  "TestClub",
		Score: 0,
	}
	for i := 0; i < 100; i++ {
		clubEvent := club.Event{
			Club:      &clubInfo,
			EventType: club.EventUserLogin,
		}
		clubServer.InChan(clubEvent)
		time.Sleep(time.Second * 1)
		if i == 10 {
			//模拟,关闭工会协程,暂停请求处理
			clubServer.CloseServer()
		} else if i == 20 {
			//模拟,停止工会协程
			clubServer.StopServer()
		} else if i == 30 {
			//模拟,重启工会协程
			clubServer.ReOpen()
		}

	}

}

func main() {
	//clubTest()
	club2Test()
}

// 创建一个ClubServer协程,并执行玩法
func club2Test() {
	//启动一个协程服务
	club := club2.Club2{}
	clubServer := club2.CreateClubServer(&club)

	//通过两个channel堵塞实现的方式
	req := club2.Request{}
	createServerRes, error := clubServer.SendRequest(req)
	if error != nil {
		log.Println("createServerRes has Error: ", error)
	} else {
		log.Println("createServerRes Result: ", createServerRes)
	}
	time.Sleep(time.Second * 10)
}
