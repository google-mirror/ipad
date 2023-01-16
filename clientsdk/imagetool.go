package clientsdk

import (
	"bytes"
	"github.com/lunny/log"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"

	"feiyu.com/wx/clientsdk/baseinfo"
	"github.com/nfnt/resize"
)

// CreateThumbImage 生成缩略图
func CreateThumbImage(fileData []byte) *baseinfo.ThumbItem {
	filetype := http.DetectContentType(fileData)
	if filetype == "image/png" {
		return CreatePngThumbImage(fileData)
	}

	if filetype == "image/jpeg" {
		return CreateJpegThumbImage(fileData)
	}

	return nil
}

// CreateJpegThumbImage 创建JPEG图片的缩略图
func CreateJpegThumbImage(srcImage []byte) *baseinfo.ThumbItem {
	// 获取图片数据
	srcImg, _ := jpeg.Decode(bytes.NewBuffer(srcImage))

	// 缩略图宽/高设置成 120
	width := uint(0)
	height := uint(120)
	if srcImg.Bounds().Dx() > srcImg.Bounds().Dy() {
		width = 120
		height = 0
	}

	// 生产缩略图
	thumbImage := resize.Resize(width, height, srcImg, resize.Lanczos3)
	emptyBuff := bytes.NewBuffer(nil)
	// write new image to file
	jpeg.Encode(emptyBuff, thumbImage, nil)

	retThumbItem := &baseinfo.ThumbItem{}
	retThumbItem.Data = emptyBuff.Bytes()
	retThumbItem.Width = int32(thumbImage.Bounds().Dx())
	retThumbItem.Height = int32(thumbImage.Bounds().Dy())

	return retThumbItem
}

// CreatePngThumbImage 创建png图片的缩略图
func CreatePngThumbImage(srcImage []byte) *baseinfo.ThumbItem {
	// 获取图片数据
	srcImg, _ := png.Decode(bytes.NewBuffer(srcImage))

	// 缩略图宽/高设置成 120
	width := uint(0)
	height := uint(120)
	if srcImg.Bounds().Dx() > srcImg.Bounds().Dy() {
		width = 120
		height = 0
	}

	// 生产缩略图
	thumbImage := resize.Resize(width, height, srcImg, resize.Lanczos3)
	emptyBuff := bytes.NewBuffer(nil)
	// write new image to file
	png.Encode(emptyBuff, thumbImage)

	// 返回数据
	retThumbItem := &baseinfo.ThumbItem{}
	retThumbItem.Data = emptyBuff.Bytes()
	retThumbItem.Width = int32(thumbImage.Bounds().Dx())
	retThumbItem.Height = int32(thumbImage.Bounds().Dy())

	return retThumbItem
}

// GetImageBounds 获取图片数据的宽高
func GetImageBounds(fileData []byte) (uint32, uint32) {
	if fileData == nil || len(fileData) <= 0 {
		log.Error("获取图片数据的宽高error!")
		return uint32(400), uint32(500)
	}
	filetype := http.DetectContentType(fileData)
	var srcImg image.Image
	if filetype == "image/png" {
		srcImg, _ = png.Decode(bytes.NewBuffer(fileData))
	}

	if filetype == "image/jpeg" {
		srcImg, _ = jpeg.Decode(bytes.NewBuffer(fileData))
	}

	if filetype == "image/gif" {
		srcImg, _ = gif.Decode(bytes.NewBuffer(fileData))
	}

	return uint32(srcImg.Bounds().Dx()), uint32(srcImg.Bounds().Dy())
}
