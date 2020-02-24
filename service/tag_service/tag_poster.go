package tag_service

import (
	"go_gin_base/pkg/file"
	"go_gin_base/pkg/qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

const (
	PosterPrefix = "poster"
)

type TagPoster struct {
	PosterName string
	*Tag
	Qr *qrcode.QrCode
}

func NewTagPoster(posterName string, tag *Tag, qr *qrcode.QrCode) *TagPoster {
	return &TagPoster{
		PosterName: posterName,
		Tag:        tag,
		Qr:         qr,
	}
}

func (t *TagPoster) CheckMergedImage(path string) bool {
	return file.CheckExist(path + t.PosterName)
}

func (t *TagPoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(t.PosterName, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type TagPosterBg struct {
	Name string
	*TagPoster
	*Rect
	*Pt
}
type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

func NewTagPosterBg(name string, tp *TagPoster, rect *Rect, pt *Pt) *TagPosterBg {
	return &TagPosterBg{
		Name:      name,
		TagPoster: tp,
		Rect:      rect,
		Pt:        pt,
	}
}

func (t *TagPosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullPath()
	fileName, path, err := t.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}
	if !t.CheckMergedImage(path) {
		mergedF, err := t.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()

		bgF, err := file.MustOpen(t.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()

		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}

		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}

		jpg := image.NewRGBA(image.Rect(t.Rect.X0, t.Rect.Y0, t.Rect.X1, t.Rect.Y1))
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(t.Pt.X, t.Pt.Y)), draw.Over)
		jpeg.Encode(mergedF, jpg, nil)
	}
	return fileName, path, nil
}
