package main

import (
	"fmt"
	"time"
	"pool"
	"sync"
)

func main(){
	fmt.Println("work pool test")
	var myPool pool.WorkPool
	myPool.Init(3,10)

	//创建task
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i:=0; i < 1000; i++{
			task := pool.NewTask(func() error{
				fmt.Printf("%v\n",time.Now().UnixNano())
				return nil
			})
			if false == myPool.AddTask(task){
				break
			}
			time.Sleep(time.Millisecond * 10)
		}
		wg.Done()
	}()

	//time.Sleep(time.Second * 3)
	//启动pool
	myPool.Run()

	//逆初始化
	myPool.Fini()

	//等待协程关闭信号
	fmt.Println("wait for sign")
	wg.Wait()
	fmt.Println("work pool test end")
	return
}

