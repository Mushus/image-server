package memory

import (
	"github.com/Mushus/image-server/server/internal"
)

type ConvertParamsRepository struct {
	dict ParamsDict
}

type ParamsDict map[string]*internal.ConvertParams

type Config interface {
	ConvertParams() ParamsDict
}

func ProvideConvertParamsRepository(config Config) (internal.ConvertParamsRepository, error) {
	return &ConvertParamsRepository{dict: config.ConvertParams()}, nil
}

func (r *ConvertParamsRepository) Get(name string) *internal.ConvertParams {
	params, ok := r.dict[name]
	if !ok {
		return internal.GetDefaultConvertParams()
	}
	return params
}
