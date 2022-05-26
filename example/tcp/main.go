// Package tcp @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/5/26 17:13
package main

import (
	"log"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Println("err:Listen", err)
		return
	}
	for {
		accept, err := listen.Accept()
		if err != nil {
			log.Println("err:Accept", err)
		}
		b := make([]byte, 1024)
		read, err := accept.Read(b)
		if err != nil {
			log.Println("err:Read", err)
		}
		log.Println(read, string(b))
		_, err = accept.Write([]byte("<html>Hello</html>"))
		if err != nil {
			log.Println("err:Read", err)
		}
	}
}
