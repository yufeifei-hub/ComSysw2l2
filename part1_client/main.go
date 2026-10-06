package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func receiveMessages(conn net.Conn) {
	// 单独goroutine：一直监听服务器发来消息
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("服务器连接断开")
			return
		}
		msg := string(buf[:n])
		fmt.Printf("\n[收到服务器消息]: %s", msg)
	}
}

func main() {
	// 主动拨号连接本机的8080服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	// 启动协程，后台接收服务器消息
	go receiveMessages(conn)

	// 新建读取器，读取键盘输入
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("已连接服务器，请输入消息，回车发送：")

	for {
		// 读取一行键盘输入（直到按回车）
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read input error:", err)
			return
		}
		// 发送这一行文字给服务器
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}

}
