package vips

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/Mushus/image-server/server/internal"
	"github.com/davidbyttow/govips/v2/vips"
)

func ProvideConverter() internal.Converter {
	vips.LoggingSettings(nil, vips.LogLevelWarning)
	vips.Startup(&vips.Config{})
	return &Converter{}
}

type Converter struct {
}

func (c *Converter) Process(image *internal.Image, params *internal.ConvertParams) (*internal.Image, error) {
	img, err := vips.NewImageFromReader(image.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot load image from reader: %w", err)
	}

	// Exif 見て回転
	if err := img.AutoRotate(); err != nil {
		return nil, fmt.Errorf("cannot auto rotate from image: %w", err)
	}

	target := calcResize(img, params)
	// 縮小
	switch params.Crop {
	case internal.CropCover:
		if err := img.Thumbnail(target.width, target.height, vips.InterestingCentre); err != nil {
			return nil, fmt.Errorf("cannot resize image: %w", err)
		}
	case internal.CropContain:
		current := size{width: img.Width(), height: img.Height()}
		resize := getContainSize(current, params)
		// 縮小
		if err := img.Thumbnail(resize.width, resize.height, vips.InterestingAll); err != nil {
			return nil, fmt.Errorf("cannot resize image: %w", err)
		}
		// 塗り足す
		left := (target.width - img.Width()) / 2
		top := (target.height - img.Height()) / 2
		if err := img.Embed(left, top, target.width, target.height, vips.ExtendBackground); err != nil {
			return nil, fmt.Errorf("cannot resize image: %w", err)
		}
	}

	b, _, err := img.ExportPng(vips.NewPngExportParams())
	if err != nil {
		return nil, fmt.Errorf("cannot export webp: %w", err)
	}

	buf := bytes.NewBuffer(b)

	return &internal.Image{Body: io.NopCloser(buf), MIME: "image/webp"}, nil
}

func (c *Converter) export(img *vips.ImageRef, params *internal.ConvertParams) (*internal.Image, error) {
	// format := img.Format()
	hasAlpha := img.HasAlpha()

	switch {
	case params.Accept.Has("image/avif"):
		b, _, err := img.ExportAvif(&vips.AvifExportParams{
			StripMetadata: true,
			Quality:       params.GetAvifQuality(),
			Lossless:      params.GetLossless(),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot export image as avif: %w", err)
		}

		return &internal.Image{
			Body: io.NopCloser(bytes.NewBuffer(b)),
			MIME: "image/avif",
		}, nil

	case params.Accept.Has("image/webp"):
		b, _, err := img.ExportWebp(&vips.WebpExportParams{
			StripMetadata:   true,
			Quality:         params.GetWebpQuality(),
			Lossless:        params.GetLossless(),
			ReductionEffort: 4,
		})

		if err != nil {
			return nil, fmt.Errorf("cannot export image as webp: %w", err)
		}

		return &internal.Image{
			Body: io.NopCloser(bytes.NewBuffer(b)),
			MIME: "image/webp",
		}, nil

	case params.Accept.Has("image/jpeg") && !hasAlpha:
		b, _, err := img.ExportJpeg(&vips.JpegExportParams{
			StripMetadata: true,
			Quality:       params.GetJpegQuality(),
			Interlace:     true,
		})

		if err != nil {
			return nil, fmt.Errorf("cannot export image as jpeg: %w", err)
		}

		return &internal.Image{
			Body: io.NopCloser(bytes.NewBuffer(b)),
			MIME: "image/jpeg",
		}, nil

	case params.Accept.Has("image/png"):
		b, _, err := img.ExportPng(&vips.PngExportParams{
			StripMetadata: true,
			Compression:   6,
			Interlace:     true,
		})

		if err != nil {
			return nil, fmt.Errorf("cannot export image as png: %w", err)
		}

		return &internal.Image{
			Body: io.NopCloser(bytes.NewBuffer(b)),
			MIME: "image/png",
		}, nil
	}

	return nil, errors.New("accept image format is not supported")
}

type size struct {
	width  int
	height int
}

// Contain時のサイズ調節を行う
func getContainSize(current size, params *internal.ConvertParams) size {
	scaleW := float64(params.GetWidth()) / float64(current.width)
	scaleH := float64(params.GetWidth()) / float64(current.height)
	if scaleW > scaleH {
		return size{
			width:  int(scaleH * float64(current.width)),
			height: int(scaleH * float64(current.height)),
		}
	}

	return size{
		width:  int(scaleW * float64(current.width)),
		height: int(scaleW * float64(current.height)),
	}
}

func calcResize(img *vips.ImageRef, params *internal.ConvertParams) size {
	iw := img.Width()
	ih := img.Height()
	s := size{
		width:  params.GetWidth(),
		height: params.GetHeight(),
	}

	if s.width == 0 && s.height == 0 {
		s.width = iw
		s.height = ih
	}

	if s.width == 0 {
		s.width = int(float64(params.GetHeight()) * float64(ih) / float64(iw))
	} else if s.height == 0 {
		s.height = int(float64(params.GetWidth()) * float64(iw) / float64(ih))
	}

	if !params.GetUpscale() {
		if s.width > iw {
			s.width = iw
		}
		if s.height > ih {
			s.height = ih
		}
	}

	return s
}
