package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	process(time.Second)
}

func process(timeout time.Duration) bool {
	ch := make(chan bool, 1)
	go func() {
		time.Sleep(timeout + time.Second)
		ch <- true
		fmt.Println("exit goroutine")
	}()

	select {
	case result := <-ch:
		return result
	case <-time.After(timeout):
		return false
	}
}

/**
有一道经典的使用 Channel 进行任务编排的题，你可以尝试做一下：有四个
goroutine，编号为 1、2、3、4。每秒钟会有一个 goroutine 打印出它自己的编号，
要求你编写一个程序，让输出的编号总是按照 1、2、3、4、1、2、3、4、……的顺序打
印出来。
*/
func TestChannel2(t *testing.T) {
	ticker := time.NewTicker(time.Second)
	for {
		for i := 1; i < 5; i++ {
			<-ticker.C
			go func(i int) {
				t.Log(i)
			}(i)
		}
	}
}

func TestChannel3(t *testing.T) {
	//ticker := time.NewTicker(time.Second)
	ch1 := make(chan int)
	for {
		for i := 0; i < 4; i++ {
			//<-ticker.C
			time.Sleep(time.Second)
			go func(d int) {
				ch1 <- d + 1
			}(i)
			t.Log(<-ch1)
		}
	}
}

type Token struct{}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch         // 取得令牌
		fmt.Println((id + 1)) // id从1开始
		time.Sleep(time.Second)
		nextCh <- token
	}
}
func TestChannel5(t *testing.T) {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}
	// 创建4个worker
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}
	//首先把令牌交给第一个worker
	chs[0] <- struct{}{}
	select {}
}

// Test that a bug tearing down a ticker has been fixed. This routine should not deadlock.
func TestTeardown(t *testing.T) {
	Delta := 100 * time.Millisecond
	if testing.Short() {
		Delta = 20 * time.Millisecond
	}
	for i := 0; i < 3; i++ {
		ticker := time.NewTicker(Delta)
		<-ticker.C
		ticker.Stop()
	}
}

func TestChannel4(t *testing.T) {
	var ch1 = make(chan int, 10)
	var ch2 = make(chan int, 10)
	// 创建SelectCase
	var cases = createCases(ch1, ch2)
	// 执行10次select
	for i := 0; i < 10; i++ {
		chosen, recv, ok := reflect.Select(cases)
		if recv.IsValid() { // recv case
			fmt.Println("recv:", cases[chosen].Dir, recv, ok)
		} else { // send case
			fmt.Println("send:", cases[chosen].Dir, ok)
		}
	}
}

func createCases(chs ...chan int) []reflect.SelectCase {
	var cases []reflect.SelectCase
	// 创建recv case
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	// 创建send case
	for i, ch := range chs {
		v := reflect.ValueOf(i)
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: v,
		})
	}
	return cases
}

//func TestChannel6(t *testing.T) {
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
