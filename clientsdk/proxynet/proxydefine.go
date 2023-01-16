package proxynet

import (
	"feiyu.com/wx/clientsdk/baseutils"
	"fmt"
	"github.com/lunny/log"
	"golang.org/x/net/proxy"
	"net/url"
)

// WXProxyInfo 代理信息
type WXProxyInfo struct {
	ProxyUrl string
	// 代理城市
	CityName string
	// 代理IP
	ProxyIP string
	// 代理端口
	ProxyPort uint32
	// 用户名
	UserName string
	// 密码
	Password string
}

func ParseWXProxyInfo(url string) *WXProxyInfo {
	return &WXProxyInfo{
		ProxyUrl: url,
	}
}

// NewWXProxyInfo 新建代理信息
func NewWXProxyInfo(cityName string, ip string, port uint32, userName string, password string) *WXProxyInfo {
	return &WXProxyInfo{
		CityName:  cityName,
		ProxyIP:   ip,
		ProxyPort: port,
		UserName:  userName,
		Password:  password,
	}
}

// GetDialer 获取代理
func (wxpi *WXProxyInfo) GetDialer() proxy.Dialer {
	if len(wxpi.ProxyUrl) > 0 {
		urlproxy, err := url.Parse(wxpi.ProxyUrl)
		if err != nil {
			return nil
		}
		dialer, err := proxy.FromURL(urlproxy, proxy.Direct)
		if err != nil {
			log.Info("获取代理error=", err.Error())
			return nil
		}
		return dialer
	}

	auth := &proxy.Auth{
		User:     wxpi.UserName,
		Password: wxpi.Password,
	}

	// 创建拨号器
	proxyAddr := fmt.Sprintf("%s:%d", wxpi.ProxyIP, wxpi.ProxyPort)
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, auth, proxy.Direct)
	if err != nil {
		baseutils.PrintLog("GetDialer err: " + err.Error())
		return nil
	}
	return dialer
}
