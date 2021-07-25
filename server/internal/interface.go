package internal

import (
	"io"
)

type Crop int

const (
	CropUnknown = iota
	CropContain
	CropCover
)

type ConvertParams struct {
	Width       *int
	Height      *int
	Crop        Crop
	Accept      MIMEs
	Quality     *int
	WebpQuality *int
	AvifQuality *int
	JpegQuality *int
	Lossless    *bool
	Upscale     *bool
}

func GetDefaultConvertParams() *ConvertParams {
	return &ConvertParams{}
}

func (p *ConvertParams) Overwrite(src ConvertParams) {
	if src.Width != nil {
		p.Width = src.Width
	}
	if src.Height != nil {
		p.Height = src.Height
	}
	if src.Crop != CropUnknown {
		p.Crop = src.Crop
	}
	if len(src.Accept) == 0 {
		p.Accept = src.Accept
	}
	if src.Quality != nil {
		p.Quality = src.Quality
	}
	if src.WebpQuality != nil {
		p.WebpQuality = src.WebpQuality
	}
	if src.AvifQuality != nil {
		p.AvifQuality = src.AvifQuality
	}
	if src.JpegQuality != nil {
		p.JpegQuality = src.JpegQuality
	}
	if src.Lossless != nil {
		p.Lossless = src.Lossless
	}
	if src.Upscale != nil {
		p.Upscale = src.Upscale
	}
}

func (p *ConvertParams) GetWidth() int {
	if p.Width == nil {
		return 0
	}
	return *p.Width
}

func (p *ConvertParams) GetHeight() int {
	if p.Height == nil {
		return 0
	}
	return *p.Height
}

func (p *ConvertParams) GetCrop() Crop {
	if p.Crop == CropUnknown {
		return CropCover
	}
	return p.Crop
}

func (p *ConvertParams) GetQuality() int {
	if p.Quality == nil {
		return 85
	}
	return *p.Quality
}

func (p *ConvertParams) GetWebpQuality() int {
	if p.WebpQuality == nil {
		return p.GetQuality()
	}
	return *p.WebpQuality
}

func (p *ConvertParams) GetAvifQuality() int {
	if p.AvifQuality == nil {
		return p.GetQuality()
	}
	return *p.AvifQuality
}

func (p *ConvertParams) GetJpegQuality() int {
	if p.JpegQuality == nil {
		return p.GetQuality()
	}
	return *p.JpegQuality
}

func (p *ConvertParams) GetLossless() bool {
	if p.Lossless == nil {
		return false
	}
	return *p.Lossless
}

func (p *ConvertParams) GetUpscale() bool {
	if p.Upscale == nil {
		return true
	}
	return *p.Upscale
}

type Image struct {
	MIME string
	Body io.ReadCloser
}

type Storage interface {
	Get(path string) (io.ReadCloser, error)
	Put(path string, file io.ReadSeeker) error
}

type Converter interface {
	Process(image *Image, params *ConvertParams) (*Image, error)
}

type ConvertParamsRepository interface {
	Get(name string) *ConvertParams
}
