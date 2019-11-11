package main

import (
	"encoding/json"
	"gochat.udp/logger"
	"gochat.udp/models"
	"gochat.udp/utils"
	"net"
	"time"
)

var (
	allIPList    map[string]struct{}
	allIDMapAddr map[int]models.Addr
	allIDMapUDPAddr   map[int]*net.UDPAddr
	msgchan      chan models.Message
	conn         *net.UDPConn
)

func init() {
	allIPList = make(map[string]struct{}, 8)
	allIDMapAddr = make(map[int]models.Addr, 8)
	allIDMapUDPAddr = make(map[int]*net.UDPAddr, 8)
	msgchan = make(chan models.Message, 8)
}

func main() {
	var err error
	addr := utils.Buildudpaddr("106.13.230.225:9713")
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
		allIDMapUDPAddr[message.From] = raddr

		if _, ok := allIPList[raddr.String()]; !ok {
			if err != nil {
				logger.Error("unmarshal addr err ", err)
				continue
			}
			allIPList[raddr.String()] = struct{}{}

			address := models.Addr{}
			address.IP = raddr
			address.ID = message.From
			address.Timestamp = time.Now().Nanosecond()
			allIDMapAddr[address.ID] = address
		}
		msgchan <- message
	}
}

func msg_center() {
	for {
		message := <-msgchan
		msgsend(message)
	}
}

func msgsend(msg models.Message) {
	ip, ok := allIDMapUDPAddr[msg.To]
	if !ok {
		logger.Error("not have user!")
	}
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("marshal err ", err)
	}
	n, err := conn.WriteToUDP(data, ip)
	if err != nil {
		logger.Errorf("write err %v write %d data", err, n)
	}
}
