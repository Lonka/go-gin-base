package upload

import (
	"go_gin_base/pkg/file"
	"go_gin_base/pkg/setting"
	"go_gin_base/pkg/util"
	"mime/multipart"
	"strings"
)

// link path
func GetImageFullUrl(name string) string {
	return setting.App.PrefixUrl + "/" + GetImagePath() + name
}

func GetImageMD5Name(name string) string {
	ext := file.GetExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

func GetImagePath() string {
	return setting.App.ImageSavePath
}

// file path
func GetImageFullPath() string {
	return setting.App.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(filename string) bool {
	ext := file.GetExt(filename)
	for _, allowExt := range setting.App.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		return false
	}
	return size <= setting.App.ImageMaxSize
}
