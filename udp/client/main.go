package main

import (
	"fmt"
	"net"
)

// UDP 客户端
func main() {
	addr, err := net.ResolveUDPAddr("udp", "115.199.111.153:2323")
	if err != nil {
		fmt.Printf("地址解析失败，err: %v", err)
		return
	}
	socket, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println("连接UDP服务器失败，err: ", err)
		return
	}
	defer socket.Close()
	sendData := []byte("Hello Server")
	_, err = socket.Write(sendData) // 发送数据
	if err != nil {
		fmt.Println("发送数据失败，err: ", err)
		return
	}

	data := make([]byte, 4096)
	n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据
	if err != nil {
		fmt.Println("接收数据失败, err: ", err)
		return
	}
	fmt.Printf("recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}
