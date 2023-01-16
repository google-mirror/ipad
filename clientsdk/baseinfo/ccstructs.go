package baseinfo

// DylibInfo DylibInfo
type DylibInfo struct {
	S string
	U string
}

// ClientCheckDataInfo ClientCheckDataInfo
type ClientCheckDataInfo struct {
	FileSafeAPI          string
	DylibSafeAPI         string
	OSVersion            string
	Model                string
	CoreCount            uint32
	VendorID             string
	ADId                 string
	NetType              uint32
	IsJaiBreak           uint32
	BundleID             string
	Device               string
	DisplayName          string
	Version              uint32
	PListVersion         uint32
	USBState             uint32
	HasSIMCard           uint32
	LanguageNum          string
	LocalCountry         string
	IsInCalling          uint32
	WechatUUID           string
	APPState             uint32
	IllegalFileList      string
	EncryptStatusOfMachO uint32
	Md5OfMachOHeader     string
	DylibInfoList        []*DylibInfo
}
