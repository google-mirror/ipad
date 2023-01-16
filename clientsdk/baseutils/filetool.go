package baseutils

import (
	"github.com/lunny/log"
	"io/ioutil"
	"os"
)

// ReadFile ReadFile
func ReadFile(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// WriteToFile WriteToFile
func WriteToFile(data []byte, fileName string) {
	err := ioutil.WriteFile(fileName, data, 0666) //写入文件(字节数组)
	if err != nil {
		log.Info("WriteToFile failed: ", err)
	}
}
