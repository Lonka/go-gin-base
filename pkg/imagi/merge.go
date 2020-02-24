package imagi

import (
	"fmt"
	"go_gin_base/pkg/file"
	"image"
	"image/draw"
	"image/jpeg"
)

func (ig Imagi) Merge(imgs []image.Image) (string, string, error) {

	for i, v := range imgs {
		if v == nil {
			return "", "", fmt.Errorf("image: %v on %d is null", v, i)
		}
	}
	bg := imgs[0]
	mergeImg := image.NewRGBA(image.Rect(0, 0, bg.Bounds().Max.X, bg.Bounds().Max.Y))
	draw.Draw(mergeImg, mergeImg.Bounds(), bg, image.ZP, draw.Src)
	for i, v := range imgs {
		if i > 0 {
			draw.Draw(mergeImg, mergeImg.Bounds(), v, v.Bounds().Min.Sub(image.Pt(ig.MergePt.X, ig.MergePt.Y)), draw.Src)
		}
	}

	mergeF, err := file.MustOpen(ig.Name, GetImagiFullPath())
	if err != nil {
		return "", "", err
	}
	defer mergeF.Close()

	err = jpeg.Encode(mergeF, mergeImg, nil)
	if err != nil {
		return "", "", err
	}
	return ig.Name, GetImagiFullPath(), nil
}
