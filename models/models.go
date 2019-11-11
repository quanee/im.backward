package models

import "net"

type Message struct {
	MsgType int    `json:"msgtype"`
	Msg     []byte `json:"msg"`
	Time    int    `json:"timestamp"`
	To      int    `json:"to"`
	From    int    `json:"from"`
}

type ChatMessage struct {
	From   int    `json:"from"`
	To     int    `json:"to"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

type Addr struct {
	ID        int    `json:"id"`
	IP        *net.UDPAddr `json:"ip"`
	Timestamp int    `json:"timestamp"`
}

type ACK struct {
	seq int
}

type KeepLive struct {
	ID int
	conn *net.UDPConn
}

type Login struct {
	UserName string `json:"username"`
	PassWord string `json:"passwd"`
}

var TIPS = "~> "
