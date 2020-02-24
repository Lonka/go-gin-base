package qrcode

import (
	"go_gin_base/pkg/file"
	"go_gin_base/pkg/setting"
	"go_gin_base/pkg/util"
	"image/jpeg"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	QrCodeEXT = ".jpg"
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    QrCodeEXT,
	}
}

func GetQrCodePath() string {
	return setting.App.QrCodeSavePath
}

//link path
func GetQrCodeFullUrl(name string) string {
	return setting.App.PrefixUrl + "/" + GetQrCodePath() + name
}

//file path
func GetQrCodeFullPath() string {
	return setting.App.RuntimeRootPath + GetQrCodePath()
}

func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.Ext
	src := path + name
	if file.CheckExist(src) == false {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}
		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}
		f, err := file.MustOpen(name, path)
		if err != nil {
			return "", "", err
		}
		defer f.Close()
		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}
	return name, path, nil
}
