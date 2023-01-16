package mmtls

import (
	"bytes"
	"github.com/lunny/log"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// MMHTTPPost mmtls方式发送数据包
func MMHTTPPost(mmInfo *MMInfo, data []byte) ([]byte, error) {
	requestURL := "http://" + mmInfo.ShortHost + mmInfo.ShortURL
	request, err := http.NewRequest("POST", requestURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("UserAgent", "MicroMessenger Client")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Cache-Control", "no-cache")
	//request.Header.Add("Connection", "close")
	request.Header.Add("content-type", "application/octet-stream")
	request.Header.Add("Upgrade", "mmtls")
	request.Header.Add("Host", mmInfo.ShortHost)
	httpTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     time.Second * 60,
	}
	//设定代理
	if mmInfo.Dialer != nil {
		log.Info("--->走代理")
		httpTransport.Dial = mmInfo.Dialer.Dial
	}

	client := &http.Client{Transport: httpTransport, Timeout: time.Second * 5}
	resp, err := client.Do(request)
	if err != nil {
		log.Debug(err.Error())
		return nil, err
	}
	// 接收响应
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 返回响应数据
	return body, nil
}

// MMHTTPPostData mmtls短链接方式发送请求数据包
func MMHTTPPostData(mmInfo *MMInfo, url string, data []byte) ([]byte, error) {
	// 创建HttpHandler
	httpHandler := &HTTPHandler{}
	httpHandler.URL = url
	httpHandler.Host = mmInfo.ShortHost
	httpHandler.MMPkg = data

	// 创建发送请求项列表
	sendItems, err := CreateSendPackItems(mmInfo, httpHandler)
	if err != nil {
		return []byte{}, err
	}

	// MMTLS-加密要发送的数据
	packData, err := MMHTTPPackData(mmInfo, sendItems)
	if err != nil {
		return []byte{}, err
	}
	// 发送数据
	respData, err := MMHTTPPost(mmInfo, packData)
	if err != nil {
		return nil, err
	}
	// 解包响应数据
	decodeData, err := MMDecodeResponseData(mmInfo, sendItems, respData)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return decodeData, nil
}

/*
*
纯Http请求
*/
func HTTPPost(mmInfo *MMInfo, cgi string, data []byte) ([]byte, error) {
	requestURL := "http://" + mmInfo.ShortURL + cgi
	request, err := http.NewRequest("POST", requestURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("UserAgent", "MicroMessenger Client")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "Keep-Alive")
	request.Header.Add("content-type", "application/octet-stream")
	// 发送请求
	httpTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     time.Second * 60,
	}
	// 如果有代理
	if mmInfo.Dialer != nil {
		httpTransport.Dial = mmInfo.Dialer.Dial
	}
	client := &http.Client{Transport: httpTransport, Timeout: time.Second * 5}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// 接收响应
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 返回响应数据
	return body, nil
}
