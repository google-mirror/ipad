package clientsdk

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"feiyu.com/wx/clientsdk/baseinfo"
)

// HTTPPost 短连接发送数据
func HTTPPost(userInfo *baseinfo.UserInfo, uri string, data []byte) ([]byte, error) {
	client := &http.Client{}

	requestURL := "http://" + userInfo.ShortHost + uri
	request, err := http.NewRequest("POST", requestURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("UserAgent", "MicroMessenger Client")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("content-type", "application/octet-stream")

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}

// ConnectCdnServer 链接Cdn服务器
func ConnectCdnServer(ipAddress string, port uint32) (*net.TCPConn, error) {
	strPort := strconv.Itoa(int(port))
	serverAddr := ipAddress + ":" + strPort
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// CDNRecvData 发送Cdn数据
func CDNRecvData(conn *net.TCPConn) []byte {
	// 写数据
	// 接收数据
	retData := make([]byte, 0)
	buffer := make([]byte, 25)
	count, err := conn.Read(buffer)
	if err != nil {
		return []byte{}
	}

	// 读取返回数据
	retData = append(retData, buffer[0:count]...)
	// 数据总长度
	totalLength := ParseCdnResponseDataLength(retData)
	currentLength := uint32(len(retData))
	for currentLength < totalLength {
		lessCount := totalLength - currentLength
		buffer := make([]byte, lessCount)
		count, err := conn.Read(buffer)
		if err != nil {
			return []byte{}
		}
		if count > 0 {
			retData = append(retData, buffer[0:count]...)
			currentLength = uint32(len(retData))
		} else {
			break
		}
	}

	return retData
}
