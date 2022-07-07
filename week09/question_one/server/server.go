package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"strings"
)

//如何解决粘包问题
//1. 告诉接收方包的“头尾”
//既然粘包的根本原因是接收方不知道发送的数据流的收尾，那么在发送的时候告诉接收方的发送方数据的大小就可以解决。
//具体做法就是在发送的TCP包前面再封装一层，就是加上数据包的长度data=str(len(data))+data,但是这样做的缺点就是编码成本较大，会增加程序的运行时间。
//
//2、固定缓冲区的大小
//发送方每次将不足缓冲区大小的数据包使用空字符填充，接收方每次接受也只接受固定长度的包长。这样做虽然可以解决粘包的问题，
//但是对于拆包的问题却不能很好的解决，同时会对网络产生额外的开销。
//
//3、特殊字符结尾
//在每个数据包的结尾加上特殊的字符，例如\n\n\t\t，这样接收方每次提取到这几个特殊的字符都会认为是一个完整的包结束了。
//可以解决粘包问题，同时对于半包(拆包)问题，也可以得到很好的解决。

//接收固定长度的数据包fix length
func handleConnByFixLen(conn net.Conn) {
	defer conn.Close()
	//接收的数据包的固定长度为38
	var buf [38]byte
	for {
		reader := bufio.NewReader(conn)
		_, err := reader.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Printf("read from conn failed, err:%v\n", err)
				break
			}

		}
		recv := string(buf[:])
		log.Printf("收到的数据：%v\n", recv)
	}
}

//根据特殊字符当作分隔符 delimiter based
func handleConnByDelimiter(conn net.Conn) {
	defer conn.Close()
	var buf [65535]byte
	for {
		reader := bufio.NewReader(conn)
		//接收的数据包的固定长度为38
		_, err := reader.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Printf("read from conn failed, err:%v\n", err)
				break
			}

		}
		recv := string(buf[:])
		//分隔符@#¥
		results := strings.Split(recv, "@#¥")
		l := len(results)
		for i := 0; i < l-1; i++ {
			res := results[i]
			log.Printf("收到的数据：%v,%d\n", res, len(res))
		}
	}
}

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 检查 atEOF 参数 和 数据包头部的四个字节是否 为 0x666666(我们定义的协议的魔数)
	if !atEOF && len(data) > 6 && binary.BigEndian.Uint32(data[:4]) == 0x666666 {
		var l int16
		// 读出 数据包中 实际数据 的长度(大小为 0 ~ 2^16)
		binary.Read(bytes.NewReader(data[4:6]), binary.BigEndian, &l)
		pl := int(l) + 6
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
			log.Printf("收到的数据：%v\n", string(scanner.Bytes()[6:]))
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
