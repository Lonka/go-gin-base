package imagi

import (
	"go_gin_base/pkg/file"
	"go_gin_base/pkg/setting"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"

	"github.com/golang/freetype"
)

type FontText struct {
	Text   string
	FontPt Pt
	Size   float64
}

func (ig Imagi) DrawText(img image.Image) (string, string, error) {
	fontSource := setting.App.RuntimeRootPath + setting.App.FontSavePath + "msjh.ttc"
	fontSourceBytes, err := ioutil.ReadFile(fontSource)
	if err != nil {
		return "", "", err
	}
	font, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return "", "", err
	}
	txtImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X, img.Bounds().Max.Y))
	draw.Draw(txtImg, txtImg.Bounds(), img, img.Bounds().Min, draw.Over)
	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(font)
	fc.SetClip(txtImg.Bounds())
	fc.SetDst(txtImg)
	fc.SetSrc(image.Black)
	for _, t := range ig.Texts {
		fc.SetFontSize(t.Size)
		_, err := fc.DrawString(t.Text, freetype.Pt(t.FontPt.X, t.FontPt.Y))
		if err != nil {
			return "", "", err
		}
	}

	txtF, err := file.MustOpen(ig.Name, GetImagiFullPath())
	if err != nil {
		return "", "", err
	}
	defer txtF.Close()

	err = jpeg.Encode(txtF, txtImg, nil)
	if err != nil {
		return "", "", err
	}
	return ig.Name, GetImagiFullPath(), nil
}
