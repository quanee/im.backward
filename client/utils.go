package main

import (
	"encoding/json"
	"gochat.udp/models"
	"net"
	"time"
)

func keepalive(udpConn *net.UDPConn) {
	ticker := time.Tick(30000 * time.Millisecond)
	message := models.Message{
		MsgType: models.KEEPLIVE,
		From:    LocalID,
	}
	for {
		select {
		case <-ticker:
			message.Time = time.Now().Nanosecond()
			data, _ := json.Marshal(message)
			udpConn.Write(data)
		}
	}
}
