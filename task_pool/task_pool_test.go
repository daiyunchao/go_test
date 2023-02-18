package task_pool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type TaskPool struct {
	limit       int         //限制协程数量
	processChan chan func() //处理业务的channel
	closeChan   chan bool   //执行关闭的channel
}

func NewTaskPool(limit int) *TaskPool {
	taskPool := &TaskPool{
		limit:       limit,
		processChan: make(chan func(), limit),
		closeChan:   make(chan bool),
	}

	return taskPool
}

func (t *TaskPool) Start(context context.Context) {
	//开启N个服务协程,去处理请求
	for i := 0; i < t.limit; i++ {
		go func(i int) {
			for {
				select {
				case f := <-t.processChan:
					//拿到请求就处理
					f()
					fmt.Printf("执行协程: %v\n", i)
				case err := <-context.Done():
					fmt.Println(err)
					fmt.Printf("结束: %v\n", i)
					return
				}
			}
		}(i)
	}
}

// Close 关闭
func (t *TaskPool) Close(cancel context.CancelFunc) {
	cancel()
}

func TestTaskPool(t *testing.T) {
	taskPool := NewTaskPool(10)
	text := context.Background()
	cancelCon, cancel := context.WithCancel(text)
	taskPool.Start(cancelCon)
	//执行1000个任务
	for i := 0; i < 1000; i++ {
		func(i int) {
			taskPool.processChan <- func() {
				fmt.Printf("%v执行任务\n", i)
			}
		}(i)
	}
	time.Sleep(3 * time.Second)
	//执行完成后关闭协程池子
	taskPool.Close(cancel)
	time.Sleep(3 * time.Second)
}
