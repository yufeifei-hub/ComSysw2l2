package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		// 等待新客户端连接
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			continue
		}
		// 拿到连接，丢进 conns channel，交给 main 的 select 处理
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	defer client.Close()
	scanner := bufio.NewScanner(client)

	for scanner.Scan() {
		line := scanner.Text()
		msg := Message{
			sender:  clientid,
			message: line + "\n", // 加上换行符，方便客户端区分消息边界
		}
		msgs <- msg
	}

	// 客户端断开 / 读取出错，走到这里
	// 注意：这里不能直接删map！map只能main协程修改
	// 👉 后面我们需要额外发一条特殊消息通知main删除这个client
	handleError(scanner.Err())

}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, err := net.Listen("tcp", *portPtr)
	handleError(err)
	if ln == nil {
		return
	}

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)

	// 用来给客户端分配ID，每来一个新用户就+1
	nextID := 1

	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients map
			// - start to asynchronously handle messages from this client
			clientID := nextID
			nextID++
			clients[clientID] = conn
			go handleClient(conn, clientID, msgs)

		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			for id, c := range clients {
				if id != msg.sender {
					_, err := c.Write([]byte(msg.message))
					if err != nil {
						handleError(err)
					}
				}
			}

		}
	}
}
