package scaling

import (
	"image"

	"golang.org/x/image/draw"
)

func LimitHeight(img image.Image, height int) image.Image {
	if img.Bounds().Dy() <= height {
		return img
	}
	return FitHeight(img, height)
}

func FitHeight(img image.Image, height int) image.Image {
	if img.Bounds().Dy() == height {
		return img
	}
	width := int(float64(img.Bounds().Dx()) * float64(height) / float64(img.Bounds().Dy()))
	return Resize(img, width, height)
}

func LimitWidth(img image.Image, width int) image.Image {
	if img.Bounds().Dx() <= width {
		return img
	}
	return FitWidth(img, width)
}

func FitWidth(img image.Image, width int) image.Image {
	if img.Bounds().Dx() == width {
		return img
	}
	height := int(float64(img.Bounds().Dy()) * float64(width) / float64(img.Bounds().Dx()))
	return Resize(img, width, height)
}

func Limit(img image.Image, maxWidth, maxHeight int) image.Image {
	widthRatio := float64(maxWidth) / float64(img.Bounds().Dx())
	heightRatio := float64(maxHeight) / float64(img.Bounds().Dy())
	if widthRatio < heightRatio {
		return LimitWidth(img, maxWidth)
	}
	return LimitHeight(img, maxHeight)
}

func Fit(img image.Image, width, height int) image.Image {
	widthRatio := float64(width) / float64(img.Bounds().Dx())
	heightRatio := float64(height) / float64(img.Bounds().Dy())
	if widthRatio < heightRatio {
		return FitWidth(img, width)
	}
	return FitHeight(img, height)
}

func Resize(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	return dst
}
