package imagi

import (
	"errors"
	"go_gin_base/pkg/file"
	"go_gin_base/pkg/setting"
	"image"
	"image/jpeg"
	"image/png"
)

type Imagi struct {
	Name    string
	MergePt Pt
	Texts   []FontText
}

type Pt struct {
	X int
	Y int
}

func GetImagiPath() string {
	return setting.App.ImagiSavePath
}

func GetImagiFullUrl(name string) string {
	return setting.App.PrefixUrl + GetImagiPath() + name
}

func GetImagiFullPath() string {
	return setting.App.RuntimeRootPath + GetImagiPath()
}

func OpenImage(fileName, filePath string) (image.Image, error) {
	errMsg := "image cannot read or decode!"
	imageF, err := file.MustOpen(fileName, filePath)
	if err != nil {
		return nil, errors.New(errMsg)
	}
	defer imageF.Close()
	var img image.Image
	ext := file.GetExt(fileName)
	if ext == ".jpg" || ext == ".jpeg" {
		img, err = jpeg.Decode(imageF)
	} else if ext == ".png" {
		img, err = png.Decode(imageF)
	} else {
		return nil, errors.New(errMsg)
	}
	if err != nil {
		return nil, errors.New(errMsg)
	}
	return img, nil
}
