package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	// 连接到本地NATS服务器，默认端口为4222
	nc, err := nats.Connect("nats://hw-sg-nono-test1.livenono.com:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// 订阅主题 "example"
	// 注册一个 订阅主题的执行器
	sub, err := nc.Subscribe("example", func(msg *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	// 发布消息到主题 "example"
	err = nc.Publish("example", []byte("Hello, NATS!"))
	if err != nil {
		log.Fatal(err)
	}

	// 等待一些时间以确保订阅者有足够的时间接收消息
	time.Sleep(500 * time.Millisecond)

	// 取消订阅
	sub.Unsubscribe()

	fmt.Println("Done")
}
