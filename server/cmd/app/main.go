package main

import (
	"log"

	"github.com/Mushus/image-server/server/adapter/env"
	"github.com/Mushus/image-server/server/adapter/http"
	"github.com/Mushus/image-server/server/adapter/memory"
	"github.com/Mushus/image-server/server/adapter/s3"
	"github.com/Mushus/image-server/server/adapter/vips"
)

func main() {
	server, err := setupServer()
	if err != nil {
		log.Fatalln(err)
	}

	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
}

func setupServer() (*http.Server, error) {
	config, err := env.ProvideConfig()
	if err != nil {
		return nil, err
	}
	storage := s3.ProvideStorage(config)
	converter := vips.ProvideConverter()
	convertParamRepository, err := memory.ProvideConvertParamsRepository(config)
	if err != nil {
		return nil, err
	}
	server := http.ProvideServer(config, storage, converter, convertParamRepository)

	return server, nil
}
