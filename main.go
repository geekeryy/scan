// Package scan @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/5/25 09:33
package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var root = &cobra.Command{
	Short: "u",
	Long:  "udp",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("err:非法参数列表")
			return
		}
		ScanUDP(args[1:]...)
	},
}

var sendData = []byte("Hello Server")

func main() {
	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func ScanUDP(address ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ch := make(chan ICMPResp)
	go ListenICMP(ctx, address, ch)

	m := make(map[string]ICMPResp)
	defer func() {
		for _, v := range address {
			if i, ok := m[v]; ok {
				log.Printf("address:%s status:%s \n", i.Address, i.Status)
			} else {
				log.Printf("address:%s status:%s \n", v, "available")
			}
		}
	}()

	for {
		select {
		case v, ok := <-ch:
			if ok {
				m[v.Address] = v
			}
		case <-ctx.Done():
			return
		}
	}
}

type ICMPResp struct {
	Address string
	Status  string
}

// ListenICMP 拦截ICMP报文
func ListenICMP(ctx context.Context, address []string, ch chan ICMPResp) {
	netAddr, err := net.ResolveIPAddr("ip4", "0.0.0.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.ListenIP("ip4:icmp", netAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for _, addr := range address {
		if err := TryUDP(addr); err != nil {
			fmt.Println(err)
			continue
		}
	}

	for {
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg, err := icmp.ParseMessage(1, buf[0:n])
		if err != nil {
			fmt.Println(err)
			return
		}

		header, err := ipv4.ParseHeader(buf[8:])
		if err != nil {
			fmt.Println(err)
			return
		}

		udp, err := ParseUDPMessage(buf[header.Len+8 : n])
		if err != nil {
			fmt.Println(err)
			return
		}

		if string(udp.Data) == string(sendData) {
			ch <- ICMPResp{
				Address: fmt.Sprintf("%s:%d", header.Dst.String(), udp.DesPort),
				Status:  ParseICMPCode(msg.Type, msg.Code),
			}
		}

		select {
		case <-ctx.Done():
			close(ch)
			fmt.Println(ctx.Err())
			return
		default:
			//fmt.Println(n, addr, msg.Type, msg.Code, msg.Checksum, string(marshal))
		}
	}
}

type UDPMessage struct {
	SrcPort  int
	DesPort  int
	Len      int
	CheckSum []byte
	Data     []byte
}

// ParseUDPMessage 解析UDP包
func ParseUDPMessage(b []byte) (*UDPMessage, error) {
	if len(b) < 8 {
		return nil, errors.New("invalid len")
	}
	m := &UDPMessage{}
	m.SrcPort = int(binary.BigEndian.Uint16(b[0:2]))
	m.DesPort = int(binary.BigEndian.Uint16(b[2:4]))
	m.Len = int(binary.BigEndian.Uint16(b[4:6]))
	m.CheckSum = b[6:8]
	m.Data = b[8:]
	return m, nil
}

// ParseICMPCode 解析ICMP类型和code
func ParseICMPCode(typeCode icmp.Type, code int) string {
	switch typeCode {
	case ipv4.ICMPTypeEchoReply:
		switch code {
		case 0:
			return "Echo Reply"
		}
	case ipv4.ICMPTypeDestinationUnreachable:
		switch code {
		case 0:
			return "Network Unreachable"
		case 1:
			return "Host Unreachable"
		case 2:
			return "Protocol Unreachable"
		case 3:
			return "Port Unreachable"
		}
	case ipv4.ICMPTypeEcho:
		return "Echo Request"
	}

	return fmt.Sprintf("未知CODE %s %d", typeCode, code)
}

// TryUDP 向目标端口发送UDP数据
func TryUDP(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("地址解析失败，err: %v", err)
	}
	socket, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("连接UDP服务器失败，err: %v", err)
	}
	defer socket.Close()

	_, err = socket.Write(sendData) // 发送数据
	if err != nil {
		return fmt.Errorf("发送数据失败，err: %v", err)
	}
	return nil
}
