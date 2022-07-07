package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	PackageLength   = 4
	HeaderLength    = 2
	ProtocolVersion = 2
	Operation       = 4
	SequenceID      = 4
)

func Enpack(d []byte) []byte {
	//PackageLength
	packageLen := make([]byte, PackageLength)
	binary.BigEndian.PutUint16(packageLen, uint16(len(d)))
	//HeaderLength
	headlerLen := make([]byte, HeaderLength)
	binary.BigEndian.PutUint16(headlerLen, uint16(100))

	//ProtocolVersion
	protocolVersion := make([]byte, ProtocolVersion)
	binary.BigEndian.PutUint16(protocolVersion, uint16(2))
	//Operation
	operation := make([]byte, Operation)
	binary.BigEndian.PutUint16(operation, uint16(1))
	//SequenceID
	sequenceID := make([]byte, SequenceID)
	binary.BigEndian.PutUint16(sequenceID, uint16(99))

	packetBuf := bytes.NewBuffer(packageLen)
	packetBuf.Write(headlerLen)
	packetBuf.Write(protocolVersion)
	packetBuf.Write(operation)
	packetBuf.Write(sequenceID)
	packetBuf.Write(d)

	return packetBuf.Bytes()
}

func clientDecoder() {

	conn, err := net.DialTimeout("tcp", "127.0.0.1:10000", time.Second*30)
	if err != nil {
		log.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	for i := 0; i < 10; i++ {
		message := "{\"ID\":\"" + strconv.Itoa(i) + "\",,\"Meta\":\"golang\",\"Content\":\"一条完整的数据\"}"
		data := Enpack([]byte(message))
		_, err = conn.Write(data)
		if err != nil {
			log.Printf("write failed , err : %v\n", err)
			break
		}

	}
}

func main() {

	clientDecoder()
}
