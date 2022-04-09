package main

import (
	"fmt"
	"net"
)

func StartFabricServer(port string) (net.Conn,error) {
	fmt.Println(">>>服务器信息：开启面向机器学习的服务器...")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil,fmt.Errorf("net.Listen error: %s",err.Error())
	}
	fmt.Printf(">>>服务器信息：开启成功\n>>>服务器信息：服务器监听在本机%s端口...\n",port)
	connection, err := listener.Accept()
	if err != nil {
		return nil,fmt.Errorf("net.Listen error: %s",err.Error())
	}
	fmt.Println(">>>服务器信息：服务器获取连接成功")
	return connection,nil
}
