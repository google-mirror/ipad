package mmtls

import (
	"github.com/lunny/log"

	"feiyu.com/wx/clientsdk/baseutils"
)

// ShowMMTLSExtensions 打印Extensions
func ShowMMTLSExtensions(extensionList []*Extension) {
	log.Println("--------------- ShowMMTLSExtensions in ---------------")
	extensionLength := len(extensionList)
	for index := 0; index < extensionLength; index++ {
		// PreSharedKeyExtensionType
		if extensionList[index].ExtensionType == PreSharedKeyExtensionType {
			tmpExtension, _ := PreSharedKeyExtensionDeSerialize(extensionList[index].ExtensionData)
			baseutils.ShowObjectValue(tmpExtension)
		}

		// ClientKeyShareType
		if extensionList[index].ExtensionType == ClientKeyShareType {
			tmpExtension, _ := ClientKeyShareExtensionDeSerialize(extensionList[index].ExtensionData)
			baseutils.ShowObjectValue(tmpExtension)
		}

		// ServerKeyShareType
		if extensionList[index].ExtensionType == ServerKeyShareType {
			tmpExtension, _ := ServerKeyShareExtensionDeSerialize(extensionList[index].ExtensionData)
			baseutils.ShowObjectValue(tmpExtension)
		}
	}
	log.Println("--------------- ShowMMTLSExtensions out ---------------")
}
