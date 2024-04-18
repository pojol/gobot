package mock

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
)

// 消息头定义
const (
	HEAD_LEN          = 7
	PACKET_LEN_SIZE   = 2
	PACKET_TYPE_SIZE  = 1
	CUSTOM_BYTES_SIZE = 2
	MSG_ID_SIZE       = 2
)

var (
	tcpByteOrder = "LittleEndian"
)

func getTcpByteOrder() binary.ByteOrder {
	if tcpByteOrder == "BigEndian" {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

func _readPackageLength(buf []byte) uint16 {
	var packetLen uint16

	br := bytes.NewReader(buf)
	binary.Read(br, getTcpByteOrder(), &packetLen)

	return packetLen
}

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

	for {
		msglenbuf := make([]byte, 2) // 使用局部变量

		_, err := io.ReadFull(conn, msglenbuf)
		if err != nil {
			break
		}

		packageLen := _readPackageLength(msglenbuf)
		msgbodyBuf := make([]byte, packageLen-2)
		_, err = io.ReadFull(conn, msgbodyBuf)
		if err != nil {
			break
		}

		var packetType uint8    // 1字节包类型
		var customBytes [2]byte // 2字节自定义字段
		var msgId uint16

		br := bytes.NewReader(msgbodyBuf)
		binary.Read(br, getTcpByteOrder(), &packetType)
		binary.Read(br, getTcpByteOrder(), &customBytes)
		binary.Read(br, getTcpByteOrder(), &msgId)

		f, _ := conn.File()
		fmt.Printf("tcp server recv fd:%v msg:%v \n", f.Fd(), msgId)

		// 处理新消息
		msgBody := msgbodyBuf[HEAD_LEN-PACKET_LEN_SIZE:]
		err = HandleMsg(conn, int(f.Fd()), msgId, msgBody)
		if err != nil {
			fmt.Println("handle msg err", msgId, err.Error())
		}
	}

	// 连接断开
	f, _ := conn.File()
	fmt.Println("client conn close =>", conn.RemoteAddr(), f.Fd())
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
	binary.Write(headerBuf, getTcpByteOrder(), uint16(7+len(msgBody)))
	binary.Write(headerBuf, getTcpByteOrder(), uint8(1))
	binary.Write(headerBuf, getTcpByteOrder(), custom)
	binary.Write(headerBuf, getTcpByteOrder(), msgId)
	binary.Write(headerBuf, getTcpByteOrder(), msgBody)

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

func StarTCPServer(byteOrder string, port string) net.Listener {
	if byteOrder != "" {
		tcpByteOrder = byteOrder
	}

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
