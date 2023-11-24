package mock

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"runtime"
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
func tcpHeaderHandle(conn *net.TCPConn) {
	//go keepalive(conn)
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Println("stack:", string(buf[:n]))
			log.Printf("recover from err: %v", err)
		}
	}()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if n == 0 || err != nil {
			continue
		}

		// 解析消息头
		var packetLen uint16
		var packetType uint8    // 1字节包类型
		var customBytes [2]byte // 2字节自定义字段
		var msgId uint16

		br := bytes.NewReader(buf[:n])
		binary.Read(br, binary.LittleEndian, &packetLen)
		binary.Read(br, binary.LittleEndian, &packetType)
		binary.Read(br, binary.LittleEndian, &customBytes)
		binary.Read(br, binary.LittleEndian, &msgId)

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

		f, _ := conn.File()
		fmt.Printf("tcp server recv fd:%v msg:%v \n", f.Fd(), msgId)

		if unfinishedMsg != nil {
			// 先处理缓存的不完整消息
			unfinishedMsg.msgBody = append(unfinishedMsg.msgBody, buf[HEAD_LEN:packetLen]...)
			err = HandleMsg(conn, int(f.Fd()), unfinishedMsg.msgId, unfinishedMsg.msgBody)
			if err != nil {
				fmt.Println("handle msg err", unfinishedMsg.msgId, err.Error())
			}
		}

		// 处理新消息
		msgBody := buf[HEAD_LEN:packetLen]
		err = HandleMsg(conn, int(f.Fd()), msgId, msgBody)

		buf = make([]byte, 1024) // 重置 buf
		if err != nil {
			fmt.Println("handle msg err", msgId, err.Error())
		}
	}

	// 请求循环结束,关闭连接
	f, _ := conn.File()
	fmt.Println("server conn closed", f.Fd())
	conn.Close()
}

// 封装写消息函数
func writeMsg(conn *net.TCPConn, msgId uint16, custom []byte, msgBody []byte) error {

	if len(custom) == 0 {
		custom = []byte{0, 0}
	}

	if len(custom) != 2 {
		return fmt.Errorf("custom bytes len must be 2")
	}

	// 构造消息头
	headerBuf := new(bytes.Buffer)
	binary.Write(headerBuf, binary.LittleEndian, uint16(7+len(msgBody)))
	binary.Write(headerBuf, binary.LittleEndian, uint8(1))
	binary.Write(headerBuf, binary.LittleEndian, custom)
	binary.Write(headerBuf, binary.LittleEndian, msgId)
	binary.Write(headerBuf, binary.LittleEndian, msgBody)

	// 发送消息头+消息体
	_, err := conn.Write(headerBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func HandleMsg(conn *net.TCPConn, fd int, msgId uint16, msgBody []byte) error {
	var err error

	if msgId == LoginGuest {
		err = tcpRouteGuestHandle(conn, fd, msgBody)
	} else if msgId == Hello {
		err = tcpHelloHandle(conn, fd, msgBody)
	} else if msgId == HeroInfo {
		err = tcpHeroInfoHandle(conn, fd, msgBody)
	} else if msgId == HeroLvup {
		err = tcpHeroLvupHandle(conn, fd, msgBody)
	}

	if err != nil {
		log.Printf("write to client error: %v", err)
		return err
	}

	return nil
}

func StarTCPServer(port string) net.Listener {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server listening on port " + port)
	go func() {
		for {
			// 接收新连接
			conn, err := ln.Accept()
			if err != nil {
				//fmt.Println("accept err", err)
				continue
			}

			// 为每个连接启动一个goroutine进行处理
			tcpconn := conn.(*net.TCPConn)

			f, _ := tcpconn.File()
			fmt.Println("new client conn =>", conn.RemoteAddr(), f.Fd())

			go tcpHeaderHandle(tcpconn)
		}
	}()

	return ln
}
