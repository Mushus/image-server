package env

import (
	"encoding/json"
	"fmt"

	"github.com/Mushus/image-server/server/adapter/memory"
	"github.com/Mushus/image-server/server/internal"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	EnvDisableExternalParams bool              `envconfig:"disable_external_params"`
	EnvConvertParams         string            `envconfig:"convert_params"`
	convertParams            memory.ParamsDict `envconfig:"-"`
	EnvBucket                string            `envconfig:"bucket"`
	EnvS3URL                 string            `envconfig:"s3_url`
}

func ProvideConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("image_server", &cfg); err != nil {
		return nil, err
	}

	convertParams, err := createConvertParams(cfg.EnvConvertParams)
	if err != nil {
		return nil, err
	}

	cfg.convertParams = convertParams

	return &cfg, nil
}

func (c *Config) DisableExternalParams() bool {
	return c.EnvDisableExternalParams
}

func (c *Config) ConvertParams() memory.ParamsDict {
	return c.convertParams
}

func (c *Config) Bucket() string {
	return c.EnvBucket
}

func (c *Config) S3URL() string {
	return c.EnvS3URL
}

type JSONConvertParams struct {
	Type        string  `json:"type"`
	Width       *int    `json:"w"`
	Height      *int    `json:"h"`
	Crop        *string `json:"crop"`
	Quality     *int    `json:"quality"`
	WebpQuality *int    `json:"webpq"`
	AvifQuality *int    `json:"avifq"`
	JpegQuality *int    `json:"jpegq"`
	Lossless    *bool   `json:"lossless"`
	Upscale     *bool   `json:"upscale"`
}

func (p *JSONConvertParams) crop() internal.Crop {
	if p.Crop == nil {
		return internal.CropUnknown
	}
	switch *p.Crop {
	case "contain":
		return internal.CropContain
	case "cover":
		return internal.CropCover
	default:
		return internal.CropUnknown
	}
}

func createConvertParams(rawParams string) (memory.ParamsDict, error) {
	dict := memory.ParamsDict{}

	if rawParams != "" {
		var params []JSONConvertParams
		if err := json.Unmarshal([]byte(rawParams), &params); err != nil {
			return nil, fmt.Errorf("cannot parse convert params config environment: %w", err)
		}

		for _, currentParams := range params {
			outParams := internal.GetDefaultConvertParams()
			outParams.Overwrite(internal.ConvertParams{
				Width:       currentParams.Width,
				Height:      currentParams.Height,
				Crop:        currentParams.crop(),
				Quality:     currentParams.Quality,
				WebpQuality: currentParams.WebpQuality,
				AvifQuality: currentParams.AvifQuality,
				JpegQuality: currentParams.JpegQuality,
				Lossless:    currentParams.Lossless,
				Upscale:     currentParams.Upscale,
			})
			dict[currentParams.Type] = outParams
		}
	}

	return dict, nil
}
