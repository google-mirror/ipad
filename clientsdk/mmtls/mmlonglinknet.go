package mmtls

import (
	"bufio"
	"errors"
	"github.com/lunny/log"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"feiyu.com/wx/clientsdk/baseutils"
)

var waitGroup sync.WaitGroup
var recvBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 128)
	},
}

// MMLongConnect 链接MMtls服务器
func MMLongConnect(mmInfo *MMInfo) error {
	if mmInfo.LONGPort == "" {
		longPort := []string{"80", "8080", "443"}
		rand.Seed(time.Now().UnixNano())
		mmInfo.LONGPort = longPort[rand.Intn(3)]
	}
	serverAddr := mmInfo.LongHost + ":" + mmInfo.LONGPort
	if mmInfo.Dialer != nil {
		log.Info("链接MMtls服务器->走代理")
		conn, err := mmInfo.Dialer.Dial("tcp4", serverAddr)
		if err != nil {
			baseutils.PrintLog("mmtls走代理error->" + err.Error())
			return err
		}
		mmInfo.Conn = conn
		mmInfo.reader = bufio.NewReader(conn)
		return nil
	}

	// 没有使用代理
	conn, err := net.DialTimeout("tcp4", serverAddr, time.Second*3)
	if err != nil {
		baseutils.PrintLog(err.Error())
		return err
	}
	mmInfo.Conn = conn
	mmInfo.reader = bufio.NewReader(conn)
	return nil
}

// MMTCPSendData MMTCPSendData 长链接发送数据
func MMTCPSendData(mmInfo *MMInfo, data []byte) error {
	// 连接服务器
	if mmInfo.Conn == nil {
		// 提前设置好长链接Host Port
		err := MMLongConnect(mmInfo)
		if err != nil {
			return err
		}
	}

	// 发送数据
	if err := mmInfo.Conn.SetWriteDeadline(time.Now().Add(time.Second * 10)); err != nil {
		return err
	}
	length, err := mmInfo.Conn.Write(data)
	// 判断是否出错
	if err != nil {
		return err
	}
	// 判断数据是否发送完毕
	if length != len(data) {
		return errors.New("MMTCPSendData err: length != len(data)")
	}

	return nil
}

func MMTCPRecv(mmInfo *MMInfo, length int) ([]byte, error) {
	var err error       // Reading error.
	var size int        // Reading size.
	var index int       // Received size.
	var buffer []byte   // Buffer object.
	var bufferWait bool // Whether buffer reading timeout set.
	var c = mmInfo.Conn

	if length > 0 {
		buffer = make([]byte, length)
	} else {
		buffer = recvBufPool.Get().([]byte)
		defer recvBufPool.Put(buffer)
	}

	var bufferSize = len(buffer)
	for {
		if length < 0 && index > 0 {
			bufferWait = true
			if err = c.SetReadDeadline(time.Now().Add(time.Microsecond * 10)); err != nil {
				return nil, err
			}
		}
		size, err = c.Read(buffer[index:])
		if size > 0 {
			index += size
			if length > 0 {
				// It reads til <length> size if <length> is specified.
				if index <= length {
					break
				}
			} else {
				if index >= bufferSize {
					// If it exceeds the buffer size, it then automatically increases its buffer size.
					buffer = append(buffer, make([]byte, bufferSize)...)
				} else {
					// It returns immediately if received size is lesser than buffer size.
					if !bufferWait {
						break
					}
				}
			}
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
				break
			}
			// Re-set the timeout when reading data.
			if bufferWait && isTimeout(err) {
				if err = c.SetReadDeadline(time.Time{}); err != nil {
					return nil, err
				}
				err = nil
				break
			}
		}

		// Just read once from buffer.
		if length == 0 {
			break
		}
	}
	return buffer[:index], err
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}

// MMTCPRecvItems 循环接收长链接数据
func MMTCPRecvItems(mmInfo *MMInfo) ([]*PackItem, error) {
	// RecodeHead *RecodeHead
	retItems := make([]*PackItem, 0)
	for i := 0; i < 4; i++ {
		packItem, err := MMTCPRecvOneItem(mmInfo)
		if err != nil {
			return nil, err
		}
		retItems = append(retItems, packItem)
	}
	return retItems, nil
}

// MMTCPRecvOneItem 接收一个完整的包
func MMTCPRecvOneItem(mmInfo *MMInfo) (*PackItem, error) {
	//log.Println(mmInfo.reader.Size())
	// 读取头部数据
	recordHeadData := make([]byte, 5)
	if _, err := io.ReadFull(mmInfo.reader, recordHeadData); err != nil {
		//log.Println("MMTCPRecvOneItem err: ", err)
		return nil, err
	}
	// 读取Content
	recordHead := RecordHeadDeSerialize(recordHeadData)
	bodyData := make([]byte, recordHead.Size)
	if _, err := io.ReadFull(mmInfo.reader, bodyData); err != nil {
		return nil, err
	}
	return &PackItem{
		RecordHead: recordHeadData,
		PackData:   bodyData,
	}, nil
}

// MMTCPSendReq 长链接发送请求
func MMTCPSendReq(mmInfo *MMInfo, seqId uint32, opCode uint32, data []byte) error {
	sendData, err := MMLongPackData(mmInfo, seqId, opCode, data)
	if err != nil {
		return err
	}

	// 发送数据
	err = MMTCPSendData(mmInfo, sendData)
	if err != nil {
		return err
	}

	return nil
}

// MMTCPRecvData 接受MMTLS信息
func MMTCPRecvData(mmInfo *MMInfo) (*LongRecvInfo, error) {
	// 接收响应数据
	recvItem, err := MMTCPRecvOneItem(mmInfo)
	if err != nil {
		//log.Info("接收mmtls长链接响应失败", err.Error())
		return nil, err
	}

	// 解包响应数据
	recvInfo, err := MMLongUnPackData(mmInfo, recvItem)
	if err != nil {
		//log.Info("解包响应数据失败", err.Error())
		return nil, err
	}

	return recvInfo, nil
}
