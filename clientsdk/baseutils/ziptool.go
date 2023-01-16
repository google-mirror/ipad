package baseutils

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"github.com/lunny/log"
	"io/ioutil"
)

// UnzipByteArray zip解压缩
func UnzipByteArray(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	retBytes, err := ioutil.ReadAll(reader)
	defer reader.Close()
	if err != nil {
		return nil, err
	}

	return retBytes, nil
}

// CompressByteArray zip压缩
func CompressByteArray(data []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(data)
	w.Close()
	return in.Bytes()
}

// DeflateZip DeflateZip
func DeflateZip(data []byte) []byte {
	// 一个缓存区压缩的内容
	buffer := bytes.NewBuffer(nil)

	// 创建一个flate.Writer
	flateWrite, err := flate.NewWriter(buffer, flate.BestCompression)
	if err != nil {
		log.Info("DeflateZip err = ", err)
		return nil
	}
	defer flateWrite.Close()
	flateWrite.Write(data)
	flateWrite.Flush()

	return buffer.Bytes()
}

// DeflateUnZip DeflateUnZip
func DeflateUnZip(data []byte) []byte {
	// 一个缓存区压缩的内容
	flateReader := flate.NewReader(bytes.NewReader(data))
	defer flateReader.Close()
	// 输出
	retBytes, err := ioutil.ReadAll(flateReader)
	if err != nil {
		log.Info("DeflateUnZip err = ", err)
		return nil
	}

	return retBytes
}
