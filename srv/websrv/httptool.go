package websrv

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

// TaskPost 短连接发送数据
func TaskPost(uri string, data []byte) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * 5}
	request, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("UserAgent", "Mozilla/5.0 (Windows; U; Windows NT 6.1; zh-CN; rv:1.9.2.6)")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
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
