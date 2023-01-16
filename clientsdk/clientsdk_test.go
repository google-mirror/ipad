package clientsdk

import (
	"testing"

	"feiyu.com/wx/db"
	"github.com/google/uuid"
	"github.com/lunny/log"
)

func TestLogin(T *testing.T) {

	// 登录
	// userInfo := qrcodeLoginDemo()
	// if userInfo == nil {
	// 	return
	// }

	// err := testGetCDNDnsInfo(userInfo)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// tmpURL := string("http://szzjwxsns.video.qq.com/102/20202/snsvideodownload?encfilekey=jEXicia3muM3GjTlk1Z3kYCU9cOibLhqREFrXaoUO8cxYzJn5L9SzicB4T655PgYJDA1D55rLdUch7DpINiaibSqhWRhibvAic3X2NaTQjRHqOibiaX1MzAQZ1QYjFiaxcTU5SfytJkCziceZ4AdfBtjiaxjn2AL97FZ7kwpc4VL4&token=AxricY7RBHdWLnyNFyR4AO84dRt8l7UGZ8vyXK6jzKugcws2flst5G7QicgUibiaIkqWBFsVmgS0WRE&idx=1&bizid=1023&dotrans=2&ef=15_0&hy=SZ")
	// fileData, err := SendCdnSnsVideoDownloadReuqest(userInfo, 1852871961, tmpURL)
	// SendCdnSnsVideoUploadReuqest(userInfo, fileData)
	// // 获取账号配置, 这个配置信息，还是要存储起来比较好
	// profileResp, err := getProfileDemo(userInfo)
	// if err != nil {
	// 	log.Println("获取账号信息失败")
	// 	return
	// }
	// baseutils.ShowObjectValue(profileResp)

	// // 发送同步消息请求
	// newSyncDemo(userInfo)

	// err = tokenLoginDemo(userInfo)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// // 初始化通讯录
	// initContactListDemo(userInfo)
}

func TestDB(T *testing.T) {
	db.InitDB()
	userInfo := NewUserInfo(uuid.New().String(), "", nil)
	db.SaveUserInfo(userInfo)
}

func TestPrint(T *testing.T) {
	tmpData := uint32(4289864217)
	log.Println(int32(tmpData))
}
