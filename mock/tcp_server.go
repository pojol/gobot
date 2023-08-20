package mock

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

// 消息头定义
const (
	HEAD_LEN = 7

	PACKET_LEN_OFFSET = 0
	PACKET_LEN_SIZE   = 2

	PACKET_TYPE_OFFSET = 2
	PACKET_TYPE_SIZE   = 1

	CUSTOM_BYTES_OFFSET = 3
	CUSTOM_BYTES_SIZE   = 2

	MSG_ID_OFFSET = 5
	MSG_ID_SIZE   = 2
)

// 缓存不完整消息的结构
type UnfinishedMessage struct {
	msgId   uint16
	msgBody []byte
}

var unfinishedMsg *UnfinishedMessage

// 处理函数,在一个新的goroutine中处理每个连接的请求
func tcpHeaderHandle(conn net.Conn) {
	//go keepalive(conn)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover from err: %v", err)
		}
	}()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if n == 0 || err != nil {
			if err != nil {
				fmt.Println("read ", n, err)
			}
			continue
		}

		// 解析消息头
		var packetLen uint16
		var packetType uint8    // 1字节包类型
		var customBytes [2]byte // 2字节自定义字段
		var msgId uint16

		br := bytes.NewReader(buf[:n])
		binary.Read(br, binary.BigEndian, &packetLen)
		binary.Read(br, binary.BigEndian, &packetType)
		binary.Read(br, binary.BigEndian, &customBytes)
		binary.Read(br, binary.BigEndian, &msgId)

		if n < int(packetLen) {
			// 消息不完整,缓存
			if unfinishedMsg == nil {
				unfinishedMsg = &UnfinishedMessage{
					msgId:   msgId,
					msgBody: buf[HEAD_LEN:n],
				}
			} else {
				unfinishedMsg.msgBody = append(unfinishedMsg.msgBody, buf[HEAD_LEN:n]...)
			}
			continue
		}

		if unfinishedMsg != nil {
			// 先处理缓存的不完整消息
			unfinishedMsg.msgBody = append(unfinishedMsg.msgBody, buf[HEAD_LEN:packetLen]...)
			err = HandleMsg(conn, unfinishedMsg.msgId, unfinishedMsg.msgBody)
			unfinishedMsg = nil

			if err != nil {
				break
			}
		}

		// 处理新消息
		msgBody := buf[HEAD_LEN:packetLen]
		err = HandleMsg(conn, msgId, msgBody)

		buf = make([]byte, 1024) // 重置 buf

		if err != nil {
			break
		}
	}

	// 请求循环结束,关闭连接
	fmt.Println("server conn closed")
	conn.Close()
}

// 封装写消息函数
func writeMsg(conn net.Conn, msgId uint16, custom []byte, msgBody []byte) error {

	if len(custom) == 0 {
		custom = []byte{0, 0}
	}

	if len(custom) != 2 {
		return fmt.Errorf("custom bytes len must be 2")
	}

	// 构造消息头
	headerBuf := new(bytes.Buffer)
	binary.Write(headerBuf, binary.BigEndian, uint16(7+len(msgBody)))
	binary.Write(headerBuf, binary.BigEndian, uint8(1))
	binary.Write(headerBuf, binary.BigEndian, custom)
	binary.Write(headerBuf, binary.BigEndian, msgId)
	binary.Write(headerBuf, binary.BigEndian, msgBody)

	// 发送消息头+消息体
	_, err := conn.Write(headerBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func HandleMsg(conn net.Conn, msgId uint16, msgBody []byte) error {
	var err error
	rand.Seed(time.Now().UnixMicro())

	fmt.Println("recv tcp header", msgId, string(msgBody))

	if msgId == 1000 || msgId == 1001 || msgId == 1002 {

		i := rand.Intn(2)
		if i == 0 {
			err = writeMsg(conn, 1001, []byte{}, []byte("joy"))
		} else {
			err = writeMsg(conn, 1002, []byte{}, []byte("ppp"))
		}

	}

	if err != nil {
		log.Printf("write to client error: %v", err)
		return err
	}

	return nil
}

func StarTCPServer(port string) net.Listener {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server listening on port 20008")
	go func() {
		for {
			// 接收新连接
			conn, err := ln.Accept()
			if err != nil {
				//fmt.Println("accept err", err)
				continue
			}

			fmt.Println("new client conn =>", conn.RemoteAddr())

			// 为每个连接启动一个goroutine进行处理
			go tcpHeaderHandle(conn)
		}
	}()

	return ln
}
