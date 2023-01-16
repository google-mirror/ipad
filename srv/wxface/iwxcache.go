package wxface

// IWXCache 缓存接口
type IWXCache interface {
	// 设置Qrcode信息
	SetQrcodeInfo(uuid string, qrAesKey []byte)
	// 获取二维码信息
	GetQrcodeInfo(uuid string) []byte
}
