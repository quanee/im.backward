package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gochat.udp/config"
	"gochat.udp/logger"
	"gochat.udp/models"
	"gochat.udp/utils"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	sendmsgchan chan models.Message
	recmsgchan  chan models.Message
	conn        *net.UDPConn
	LocalID     int
)

func init() {
	sendmsgchan = make(chan models.Message, 8)
	recmsgchan = make(chan models.Message, 8)
}

func main() {
	var err error
	addr := utils.Buildudpaddr(config.GetKey("host") + ":" + config.GetKey("port"))
	conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		logger.Error("dial udp err ", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("input you id: ")
	localid, err := reader.ReadString('\n')
	if err != nil {
		logger.Error("reader bufio err ", err)
	}
	LocalID, err = strconv.Atoi(strings.Trim(localid, "\r\n"))
	message := models.Message{
		Time: time.Now().Nanosecond(),
		From: LocalID,
	}
	data, err := json.Marshal(message)
	if err != nil {
		logger.Error("unmarshal message err ", err)
	}
	n, err := conn.Write(data)
	if err != nil {
		logger.Errorf("write err %v write %d data", err, n)
	}

	go keepalive(conn)
	go msg_center()
	go read_data()
	go msg_rec()
	for {
		printTips()
		localid, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("reader bufio err ", err)
			continue
		}
		msg := strings.Trim(localid, "\r\n")
		msgs := strings.Split(msg, ":")
		id, err := strconv.Atoi(strings.Trim(msgs[0], "\r\n"))
		if err != nil {
			logger.Error("ascii to int err ", err)
			continue
		}
		sendmessage := models.Message{
			Time:    time.Now().Nanosecond(),
			From:    LocalID,
			To:      id,
			MsgType: models.PRIVATEMSG,
			Msg:     []byte(msgs[1]),
		}
		sendmsgchan <- sendmessage
	}
}

func read_data() {
	buf := make([]byte, 1024)
	message := models.Message{}
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Error("read data err ", err)
		}
		err = json.Unmarshal(buf[:n], &message)
		if err != nil {
			logger.Error("unmarshal message err ", err)
			continue
		}
		switch message.MsgType {
		case models.PRIVATEMSG:
			recmsgchan <- message
		}
	}
}

func msg_rec() {
	for {
		msg := <-recmsgchan
		msg_print(msg)
	}
}

func msg_print(message models.Message) {
	fmt.Printf("\n%d: %s\n", message.From, message.Msg)
	printTips()
}

func msg_center() {
	for {
		msg := <-sendmsgchan
		msg_send(msg)
	}
}

func msg_send(message models.Message) {
	data, err := json.Marshal(message)
	if err != nil {
		logger.Error("unmarshal queryaddr err ", err)
	}
	n, err := conn.Write(data)
	if err != nil {
		logger.Errorf("write err %v write %d data", err, n)
	}
}

func printTips() {
	fmt.Printf("~-> ")
}
