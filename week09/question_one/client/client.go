package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

//固定长度
func clientFixLen() {
	data := []byte("[这里才是一个完整的数据包]")
	log.Println("package len", len(data))
	conn, err := net.DialTimeout("tcp", "127.0.0.1:10000", time.Second*30)
	if err != nil {
		log.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	for i := 0; i < 1000; i++ {
		_, err = conn.Write(data)
		if err != nil {
			log.Printf("write failed , err : %v\n", err)
			break
		}

	}
}

//分隔符
func clientDelimiter() {
	//分隔符@#¥
	data := []byte("[这里才是一个完整的数据包]@#¥")
	log.Println("package len", len(data))
	conn, err := net.DialTimeout("tcp", "127.0.0.1:10000", time.Second*30)
	if err != nil {
		log.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	for i := 0; i < 1000; i++ {
		_, err = conn.Write(data)
		if err != nil {
			log.Printf("write failed , err : %v\n", err)
			break
		}

	}
}

func clientDecoder() {
	data := []byte("[这里才是一个完整的数据包]")
	log.Println("package len", len(data))

	magic := make([]byte, 4)
	binary.BigEndian.PutUint32(magic, 0x666666)
	lenNum := make([]byte, 2)
	binary.BigEndian.PutUint16(lenNum, uint16(len(data)))
	packetBuf := bytes.NewBuffer(magic)
	packetBuf.Write(lenNum)
	packetBuf.Write(data)

	conn, err := net.DialTimeout("tcp", "127.0.0.1:10000", time.Second*30)
	if err != nil {
		log.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	for i := 0; i < 1; i++ {
		_, err = conn.Write(packetBuf.Bytes())
		if err != nil {
			log.Printf("write failed , err : %v\n", err)
			break
		}

	}
}

func main() {

	//clientFixLen()
	//clientDelimiter()
	clientDecoder()
}
