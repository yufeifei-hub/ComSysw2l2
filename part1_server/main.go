package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var clients = make([]net.Conn, 0)
var lock sync.Mutex

func handleClient(conn net.Conn) {
	defer conn.Close()
	// 新客户端加入列表
	lock.Lock()
	clients = append(clients, conn)
	lock.Unlock()
	fmt.Println("收到一个客户端连接")

	buf := make([]byte, 1024)
	for { // 循环：持续读这个客户端的数据
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("客户端断开或读取出错:", err)
			// ========== 新增：删除断开的客户端 ==========
			lock.Lock()
			newClients := make([]net.Conn, 0)
			for _, c := range clients {
				if c != conn {
					newClients = append(newClients, c)
				}
			}
			clients = newClients
			lock.Unlock()

			return
		}
		msg := string(buf[:n])
		fmt.Printf("收到客户端消息: %s", msg)
	}

}

func broadcast() {
	// 服务器主线程：读取键盘输入，广播给所有client
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		lock.Lock()
		for _, c := range clients {
			_, err := c.Write([]byte(msg))
			if err != nil {
				fmt.Println("发送给客户端失败")
			}
		}
		lock.Unlock()
	}
}

func main() {
	// 在本机8080端口开启TCP监听
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server正在监听 127.0.0.1:8080，等待客户端连接...")

	// 单独goroutine，负责服务器控制台广播消息
	go broadcast()

	// 无限循环，持续接收新客户端
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		// 开goroutine，并发处理这个客户端
		go handleClient(conn)
	}

}
