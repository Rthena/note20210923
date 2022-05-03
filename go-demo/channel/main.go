package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//func main() {
//	var ch = make(chan int, 10)
//	for i := 0; i < 10; i++ {
//		select {
//		case ch <- i:
//			fmt.Println("i=", i)
//		case v := <-ch:
//			fmt.Println(v)
//		}
//	}
//}

//func main() {
//	go func() {
//		//...... // 执行业务处理
//	}()
//	// 处理CTRL+C等中断信号
//	termChan := make(chan os.Signal)
//	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
//	<-termChan
//	// 执行退出之前的清理动作
//	doCleanup()
//	fmt.Println("优雅退出")
//}
//
//func doCleanup() {
//	fmt.Println("do clean up")
//}

func main1() {
	var closing = make(chan struct{})
	var closed = make(chan struct{})
	go func() {
		// 模拟业务处理
		for {
			select {
			case <-closing:
				return
			default:
				// ....... 业务计算
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	// 处理CTRL+C等中断信号
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	close(closing)
	// 执行退出之前的清理动作
	go doCleanup(closed)
	select {
	case <-closed:
	case <-time.After(time.Second):
		fmt.Println("清理超时，不等了")
	}
	fmt.Println("优雅退出")
}

func doCleanup(closed chan struct{}) {
	time.Sleep((time.Minute))
	close(closed)
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2: // 2个也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: //超过两个，二分法递归处理
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}
		}
	}()
	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()
	<-or(
		sig(10*time.Second),
		sig(20*time.Second),
		sig(30*time.Second),
		sig(40*time.Second),
		sig(50*time.Second),
		sig(01*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
