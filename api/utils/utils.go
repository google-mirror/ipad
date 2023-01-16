package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"fmt"
	"github.com/lunny/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

func ScanIntoGrouppost(URL string, deviceType int, usrInfo *baseinfo.UserInfo) (string, error) {
	var err error
	postValue := url.Values{
		"forBlackberry": {"forceToUsePost"},
	}
	req, err := http.PostForm(URL, postValue)
	if err != nil {
		return "", err
	}
	ua := "Mozilla/5.0 (Android8.1.0) AppleWebKit/537. 36 (KHTML, like Gecko) Chrome/41. 0.2225.0 Safari/537. 36"
	if deviceType == 1 {
		ua = "User-Agent: Mozilla/5.0 (iPad; CPU OS 12_4_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/" + baseinfo.DeviceVersionNumber + " NetType/WIFI " + usrInfo.DeviceInfo.Language + "/" + usrInfo.DeviceInfo.RealCountry
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", URL)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", ua)
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

var (
	DeviceLenError = errors.New("设备长度错误!")
)

func GenerateWxDat(device string) (string, error) {
	if len(device) < 32 {
		return "", DeviceLenError
	}
	datByte := []byte(device)
	datHex := hex.EncodeToString(datByte)
	str := "62706c6973743030d4010203040506090a582476657273696f6e58246f626a65637473592461726368697665725424746f7012000186a0a2070855246e756c6c5f1020" + datHex + "5f100f4e534b657965644172636869766572d10b0c54726f6f74800108111a232d32373a406375787d0000000000000101000000000000000d0000000000000000000000000000007f"
	return str, nil
}

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 识别手机号码
func IsMobile(mobile string) bool {
	result, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, mobile)
	if result {
		return true
	} else {
		return false
	}
}

/*
**
获取微信html页面
*/
func GetHTML(url string, headList []*wechat.GetA8KeyRespHeader, usrInfo *baseinfo.UserInfo) (string, []string, error) {
	if len(url) < 1 {
		return "", nil, errors.New("url 获取错误")
	}
	ua := "Mozilla/5.0 (Android8.1.0) AppleWebKit/537. 36 (KHTML, like Gecko) Chrome/41. 0.2225.0 Safari/537. 36"
	Language := "zh-cn,zh"
	if usrInfo.DeviceInfo != nil {
		ua = "User-Agent: Mozilla/5.0 (iPad; CPU OS 12_4_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/" + baseinfo.DeviceVersionNumber + " NetType/WIFI " + usrInfo.DeviceInfo.Language + "/" + usrInfo.DeviceInfo.RealCountry
		Language = usrInfo.DeviceInfo.Language + "," + usrInfo.DeviceInfo.RealCountry
	}
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Charset", "utf-8;q=0.7,*;q=0.3")
	//reqest.Header.Set("Accept-Encoding", "gzip, default")//这个有乱码，估计是没有解密，或解压缩
	reqest.Header.Set("Accept-Encoding", "utf-8") //这就没有乱码了
	reqest.Header.Set("Accept-Language", Language+";q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Set("Cache-Control", "max-age=0")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Host", url)
	reqest.Header.Set("User-Agent", ua)
	for _, head := range headList {
		reqest.Header.Add(head.GetName(), head.GetKey())
	}
	response, err := client.Do(reqest)
	if err == nil {
		if response != nil && response.StatusCode == 200 {
			body, _ := ioutil.ReadAll(response.Body)
			bodystr := string(body)
			v := response.Header
			b := v["Set-Cookie"]
			return bodystr, b, nil //response.Header..Get("Set-Cookie")
		}
		log.Warn("获取html内容错误:", response)
		return "", nil, errors.New("html 返回错误")
	}
	log.Warn("获取html错误:", err.Error())
	return "", nil, err
}

func HttpPost(urls string, data url.Values, cookie []string, usrInfo *baseinfo.UserInfo) string {
	client := &http.Client{}
	retest, err := http.NewRequest("POST", urls, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("Http Post NewRequest出错!")
		return ""
	}
	ua := "Mozilla/5.0 (Android8.1.0) AppleWebKit/537. 36 (KHTML, like Gecko) Chrome/41. 0.2225.0 Safari/537. 36"
	if usrInfo.DeviceInfo != nil {
		ua = "User-Agent: Mozilla/5.0 (iPad; CPU OS 12_4_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/" + baseinfo.DeviceVersionNumber + " NetType/WIFI " + usrInfo.DeviceInfo.Language + "/" + usrInfo.DeviceInfo.RealCountry
	}
	str := strings.Replace(strings.Trim(fmt.Sprint(cookie), "[]"), "HttpOnly", "", -1)
	retest.Header.Add("Cookie", str)
	retest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	retest.Header.Add("User-Agent", ua)
	resp, err := client.Do(retest)
	if err != nil {
		log.Error("Http Post Do出错!")
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Http Post请求出错!")
	}
	return string(body)
}

func BytesString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
