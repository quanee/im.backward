package main

import (
	"encoding/json"
	"gochat.udp/config"
	db "gochat.udp/database"
	"gochat.udp/logger"
	"gochat.udp/models"
	"gochat.udp/utils"
	"net"
)

var (
	allIPList    map[string]struct{}
	msgchan      chan models.Message
	conn         *net.UDPConn
)

func init() {
	allIPList = make(map[string]struct{}, 8)
	msgchan = make(chan models.Message, 8)
}

func main() {
	var err error
	addr := utils.Buildudpaddr(config.GetKey("host")+":"+config.GetKey("port"))
	conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		logger.Info("Failed to Listen")
		return
	}

	go msg_center()

	for {
		buf := make([]byte, 1024)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Error("Failed to read udp ", err)
			continue
		}
		if 0 == n {
			logger.Error("Read udp data len 0 ", err)
			continue
		}
		message := models.Message{}
		err = json.Unmarshal(buf[:n], &message)
		if err != nil {
			logger.Error("unmarshal message err ", err)
			continue
		}
		switch message.MsgType {
		case models.KEEPLIVE:
			db.KeepOnline(message.From, raddr.String())
		default:
			db.IdtoUdpAddr(message.From, raddr.String())
			db.KeepOnline(message.From, raddr.String())
			msgchan <- message
		}

	}
}

func msg_center() {
	for {
		message := <-msgchan
		msgsend(message)
	}
}

func msgsend(msg models.Message) {
	strudpaddr, ok := db.GetUdpAddrById(msg.To)
	if !ok {
		logger.Error("not have user!")
	}
	ip := utils.Buildudpaddr(strudpaddr)

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("marshal err ", err)
	}

	n, err := conn.WriteToUDP(data, ip)
	if err != nil {
		logger.Errorf("write err %v write %d data", err, n)
	}
}
