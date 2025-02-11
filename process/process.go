package process

import (
	"fmt"

	"github.com/anlakii/wallify/config"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type ImageProcessor struct {
	Config *config.Config
}

func (i *ImageProcessor) Process() error {
	imagick.Initialize()
	defer imagick.Terminate()

	coverSmallHeight := uint(float32(i.Config.Width) / 5)
	offsetX := int(float32(0.055) * float32(i.Config.Width))

	cover := imagick.NewMagickWand()
	coverSmall := cover.Clone()

	err := cover.ReadImage(i.Config.CoverPath)
	if err != nil {
		return err
	}

	err = coverSmall.ReadImage(i.Config.CoverPath)
	if err != nil {
		return err
	}

	err = cover.BlurImage(10, 5)
	if err != nil {
		return err
	}

	err = cover.ResizeImage(i.Config.Width, i.Config.Width, 0)
	if err != nil {
		return err
	}

	err = cover.CropImage(
		i.Config.Width,
		i.Config.Height,
		0,
		int(i.Config.Width-i.Config.Height)/2,
	)
	if err != nil {
		return err
	}

	err = coverSmall.ResizeImage(coverSmallHeight, coverSmallHeight, imagick.FILTER_CUBIC)
	if err != nil {
		return err
	}

	err = i.addShadow(&coverSmall, 1.015, 60, 8.75, 10, 10)
	if err != nil {
		return err
	}

	err = i.dimImage(&cover, 0.3)
	if err != nil {
		return err
	}

	err = cover.CompositeImage(coverSmall, imagick.COMPOSITE_OP_OVER, true, offsetX, int(i.Config.Height/2-(coverSmallHeight/2)))
	if err != nil {
		return err
	}

	err = cover.WriteImage(i.Config.SavePath)

	return err
}

func (i* ImageProcessor) dimImage(mw **imagick.MagickWand, dimPercent float32) error {
	overlay := imagick.NewMagickWand()
	height := (*mw).GetImageHeight()
	width := (*mw).GetImageWidth()
	pw := imagick.NewPixelWand()
	pw.SetColor(fmt.Sprintf("rgba(0, 0, 0, %f)", dimPercent))
	err := overlay.NewImage(uint(float32(width)), uint(float32(height)), pw)
	if err != nil {
		return err
	}

	err = overlay.SetImageBackgroundColor(pw)
	if err != nil {
		return err
	}

	err = overlay.CompositeImage(*mw, imagick.COMPOSITE_OP_DARKEN, true, 0, 0)
	*mw = overlay
	return err

}

func (i *ImageProcessor) addShadow(mw **imagick.MagickWand, shadowSize float32, opacity, sigma float64, x, y int) error {
	shadow := imagick.NewMagickWand()
	pw := imagick.NewPixelWand()
	pw.SetColor("black")

	height := (*mw).GetImageHeight()
	width := (*mw).GetImageWidth()
	shadowHeight := uint(float32((*mw).GetImageHeight()) * shadowSize)
	shadowWidth := uint(float32((*mw).GetImageWidth()) * shadowSize)
	err := shadow.NewImage(uint(float32(width)*shadowSize), uint(float32(height)*shadowSize), pw)
	if err != nil {
		return err
	}

	err = shadow.SetImageBackgroundColor(pw)
	if err != nil {
		return err
	}

	err = shadow.ShadowImage(opacity, sigma, x, y)
	if err != nil {
		return err
	}

	err = shadow.CompositeImage(*mw, imagick.COMPOSITE_OP_OVER, true, int((shadowWidth-width) / 2), int((shadowHeight-height) / 2))
	*mw = shadow
	return err

}
