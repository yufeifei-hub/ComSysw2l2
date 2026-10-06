package main

import (
	"bufio"
	"flag"
	"net"
	"os"
)

func read(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		println(line)
	}
	// 服务器断开，退出
}

func write(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		// 加上换行符，server那边用scanner按行读取
		_, err := conn.Write([]byte(line + "\n"))
		if err != nil {
			return
		}
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	//TODO Try to connect to the server
	// 连接服务器
	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//TODO Start asynchronously reading and displaying messages
	// 异步启动 read 协程：接收别人消息
	go read(conn)

	//TODO Start getting and sending user messages.
	// 主协程跑 write：自己打字发送
	write(conn)
}
