package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

const (
	PackageLength   = 4
	HeaderLength    = 2
	ProtocolVersion = 2
	Operation       = 4
	SequenceID      = 4
)

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 检查 atEOF 参数 )
	total := PackageLength + HeaderLength + ProtocolVersion + Operation + SequenceID
	if !atEOF && len(data) > total {
		var l int16
		// 读出 数据包中 实际数据 的长度(大小为 0 ~ 2^16)
		binary.Read(bytes.NewReader(data[:total]), binary.BigEndian, &l)
		pl := int(l) + total
		if pl <= len(data) {
			return pl, data[:pl], nil
		}
	}
	return
}

//解码器 length field based frame decoder
func handleConnByDecoder(conn net.Conn) {
	defer conn.Close()
	var buf [65535]byte
	result := bytes.NewBuffer(nil)
	for {
		reader := bufio.NewReader(conn)
		//接收的数据包的固定长度为38
		n, err := reader.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Printf("read from conn failed, err:%v\n", err)
				break
			}

		}
		result.Write(buf[0:n])
		scanner := bufio.NewScanner(result)
		scanner.Split(splitFunc)
		for scanner.Scan() {
			var packageLen, headlerLen, protocolVersion, operation, sequenceID int16
			binary.Read(bytes.NewReader(scanner.Bytes()[:4]), binary.BigEndian, &packageLen)
			binary.Read(bytes.NewReader(scanner.Bytes()[4:6]), binary.BigEndian, &headlerLen)
			binary.Read(bytes.NewReader(scanner.Bytes()[6:8]), binary.BigEndian, &protocolVersion)
			binary.Read(bytes.NewReader(scanner.Bytes()[8:12]), binary.BigEndian, &operation)
			binary.Read(bytes.NewReader(scanner.Bytes()[12:16]), binary.BigEndian, &sequenceID)
			message := string(scanner.Bytes()[16:])
			log.Printf("收到的数据：packageLen:%v,headlerLen:%v,protocolVersion:%v,operation:%v,sequenceID:%v, message:%v\n", packageLen, headlerLen, protocolVersion, operation, sequenceID, message)
		}

		result.Reset()
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		// 开始goroutine监听连接
		go handleConnByDecoder(conn)
	}

}
