package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

/************************ 交互协议 ****************************/
/* 数据包格式：包头+数据长度(4bytes)+数据 */
const (
	ConstHeader         = "www.love-program.com"  	//包头
	ConstHeaderLength   = 20			//包头长度
	ConstSaveDataLength = 4				//数据长度
)

// 封包
// header(20 bytes) + messageLength(4 bytes) + message(n bytes)
func TCPPacket(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

// 解包
func TCPUnpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstSaveDataLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
			readerChannel <- data
			i += ConstHeaderLength + ConstSaveDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

func TCPCheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}


/* TCP 同步阻塞发包*/
func TCPPacketSend(conn net.Conn, data []byte) (n int, err error) {
	return conn.Write(TCPPacket(data))
}

/* TCP 异步持续收包，客户端退出则关闭*/
func TCPPacketReceive(conn net.Conn, readerChan chan []byte) {
	defer conn.Close()

	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), "Connection read data error: ", err)
			break
		}

		tmpBuffer = TCPUnpack(append(tmpBuffer, buffer[:n]...), readerChan)
	}
}

/************************ 服务器 ****************************/
var TCPServerConns []net.Conn
func TCPServerRun() {
	TCPServerConns = make([]net.Conn, 0)
	go func() {
		for {
			text := InputString()
			JJKPrintln("count:", len(TCPServerConns))
			for index, conn := range TCPServerConns {
				_, err := TCPPacketSend(conn, []byte(text))
				if err != nil {
					JJKPrintln(err)
					TCPServerConns = append(TCPServerConns[:index], TCPServerConns[index+1:]...)
				}
			}
			JJKPrintln("broadcast ok!")
		}
	}()

	netListen, err := net.Listen("tcp", "0.0.0.0:6080")
	TCPCheckError(err)
	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		JJKPrintln(fmt.Sprintf("%v tcp connect success", conn.RemoteAddr()))
		go TCPServerHandleConnection(conn)

		TCPServerConns = append(TCPServerConns, conn)

	}
}

func TCPServerHandleConnection(conn net.Conn) {
	defer conn.Close()

	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//接收解包
	readerChan := make(chan []byte, 16)
	go TCPServerReader(readerChan)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), "Connection read data error: ", err)
			break
		}

		tmpBuffer = TCPUnpack(append(tmpBuffer, buffer[:n]...), readerChan)
	}

	JJKPrintln(fmt.Sprintf("%v client disconnect", conn.RemoteAddr()))
}

func TCPServerReader(readerChan chan []byte) {
	for {
		select {
		case data := <-readerChan:
			JJKPrintln(fmt.Sprintf("%v",string(data)))
		}
	}
}

/************************ 客户端 ****************************/

func TCPClientSend(conn net.Conn) {
	defer conn.Close()
	for {
		JJKPrintln("\n[1]quit\n[2]send json data\n[3]custom text")
		code := InputIntWithMessage("Enter Code:")
		quit := false
		switch code {
		default:
			quit = true
		case 2:
			for i := 0; i < 2; i++ {
				session := TCPClientGetSession()
				words := "{\"ID\":" + strconv.Itoa(i) + "\",\"Session\":" + session + "2015073109532345\"}"
				fmt.Println(words)
				conn.Write(TCPPacket([]byte(words)))
			}
		case 3:
			conn.Write(TCPPacket([]byte(InputStringWithMessage("Enter Message:"))))
		}
		if quit {
			break
		}
		JJKPrintln("send success!")
	}

	JJKPrintln("client exit!")
}

func TCPClientGetSession() string {
	gs1 := time.Now().Unix()
	gs2 := strconv.FormatInt(gs1, 10)
	return gs2
}

func TCPClientRun() {
	server := "0.0.0.0:6080"
	//server := "10.11.248.91:6080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")

	//接受包
	readerChan := make(chan []byte, 16)
	go TCPPacketReceive(conn, readerChan)
	go TCPClientReadChan(conn, readerChan)

	//发送包
	TCPClientSend(conn)
}

func TCPClientReadChan(conn net.Conn,readChan chan []byte)  {
	for {
		select {
		case data := <- readChan:
			JJKPrintln(string(data))
		}
	}
}