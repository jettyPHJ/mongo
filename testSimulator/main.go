package main

import (
	"context"
	"flag"
	"fmt"
	l "testSimulator/log"
	"testSimulator/simulator"
	"time"
)

// 定义命令行参数对应的变量，这三个变量都是指针类型
var clinum = flag.Int("num", 1, "Input clinum")
var frequency = flag.Int("frequency", 1, "Input frequency")
var duration = flag.Int("duration", 5, "Input duration")

func main() {
	// 把用户传递的命令行参数解析为对应变量的值
	flag.Parse()
	clients := simulator.NewSimulators(*clinum, *frequency)
	ctx, cancel := context.WithCancel(context.Background())
	go simulator.PrintGoroutinNum()
	go clients.Send(ctx)
	go CountDown(ctx, *duration)
	time.Sleep(time.Second * time.Duration(*duration))
	println("cancel......")
	cancel()
	println("finish")
}

func CountDown(ctx context.Context, secends int) {
	defer func() {
		if err := recover(); err != nil {
			l.ErrorLogger.Println("协程出错：", err)
		}
	}()
	i := 1
	for {
		i++
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			fmt.Printf("\r 剩余时间：%d", secends-i)
		}
	}
}
