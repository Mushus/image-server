package http

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Mushus/image-server/server/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Server struct {
	app *fiber.App
}

func ProvideServer(
	config Config,
	storage internal.Storage,
	converter internal.Converter,
	convertParamRepository internal.ConvertParamsRepository,
) *Server {
	disableExternalParams := config.DisableExternalParams()
	handler := Handler{
		disableExternalParams:  disableExternalParams,
		storage:                storage,
		converter:              converter,
		convertParamRepository: convertParamRepository,
	}

	app := fiber.New(fiber.Config{
		BodyLimit:             defaultMaxMemory,
		StreamRequestBody:     true,
		DisableStartupMessage: true,
	})

	app.Use(requestid.New())
	app.Use(errorResponser)
	// app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/*", handler.Get)
	app.Put("/*", handler.Put)

	return &Server{
		app: app,
	}
}

func (s *Server) Start() error {
	return s.app.Listen(":8080")
}

type Config interface {
	DisableExternalParams() bool
}

const HeaderAccept = "Accept"

const defaultMaxMemory = 24 << 20 // 32MB

type Message struct {
	Message string `json:"message"`
}

var MessageBadRequest = Message{Message: "bad request"}
var MessageInternalServerError = Message{Message: "internal server error"}
var MessageNotFound = Message{Message: "not found"}

type ConvParam struct {
	Type        string `query:"type"`
	Width       *int   `query:"w"`
	Height      *int   `query:"h"`
	Crop        string `query:"crop"`
	Quality     *int   `query:"quality"`
	WebpQuality *int   `query:"webpq"`
	AvifQuality *int   `query:"avifq"`
	JpegQuality *int   `query:"jpegq"`
	Lossless    *bool  `query:"lossless"`
	Upscale     *bool  `query:"upscale"`
}

func (c *ConvParam) GetCrop() internal.Crop {
	switch c.Crop {
	case "contain":
		return internal.CropContain
	case "cover":
		return internal.CropCover
	default:
		return internal.CropCover
	}
}

type Handler struct {
	disableExternalParams  bool
	storage                internal.Storage
	converter              internal.Converter
	convertParamRepository internal.ConvertParamsRepository
}

func (h *Handler) Put(c *fiber.Ctx) error {
	mf, err := c.MultipartForm()
	if err != nil {
		return BadRequest(c)
	}

	path := c.Params("*")
	if path == "" {
		return BadRequest(c)
	}

	imageFiles, ok := mf.File["image"]
	if !ok {
		return BadRequest(c)
	}
	// ok帰ってきてるのできっと0はいる
	img := imageFiles[0]

	if img.Filename == "" {
		return BadRequest(c)
	}

	f, err := img.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	if err := h.storage.Put(path, f); err != nil {
		return err
	}

	return c.JSON(Message{Message: "ok"})
}

func (h *Handler) Get(c *fiber.Ctx) error {
	path := c.Params("*")
	if path == "" {
		return BadRequest(c)
	}

	mimes := getAcceptMIME(c)

	var qp ConvParam
	if err := c.QueryParser(&qp); err != nil {
		return BadRequest(c)
	}

	params := h.convertParamRepository.Get(qp.Type)
	if h.disableExternalParams {
		params.Overwrite(internal.ConvertParams{
			Accept: mimes,
		})
	} else {
		params.Overwrite(internal.ConvertParams{
			Width:       qp.Width,
			Height:      qp.Height,
			Crop:        qp.GetCrop(),
			Accept:      mimes,
			Quality:     qp.Quality,
			WebpQuality: qp.WebpQuality,
			AvifQuality: qp.AvifQuality,
			JpegQuality: qp.JpegQuality,
			Lossless:    qp.Lossless,
			Upscale:     qp.Upscale,
		})
	}

	image, err := h.storage.Get(path)
	if err != nil {
		return err
	}
	defer image.Close()

	in := &internal.Image{
		Body: image,
	}

	img, err := h.converter.Process(in, params)
	if err != nil {
		return err
	}
	defer img.Body.Close()

	c.Status(200)
	c.Response().Header.SetContentType(img.MIME)
	if _, err := io.Copy(c, img.Body); err != nil {
		return err
	}

	return nil
}

func BadRequest(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(MessageBadRequest)
}

func errorResponser(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			return c.Status(http.StatusNotFound).JSON(MessageNotFound)
		}
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(MessageInternalServerError)
	}
	return nil
}

func getAcceptMIME(c *fiber.Ctx) internal.MIMEs {
	accept := c.Get(HeaderAccept)
	if accept == "" {
		return internal.MIMEs{"image/png", "image/jpeg", "image/gif"}
	}

	splitedAccepts := strings.Split(accept, ",")
	mimes := make(internal.MIMEs, 0, len(splitedAccepts))
	for _, accept := range splitedAccepts {
		mime := accept
		endOfMIME := strings.Index(mime, ";")
		if endOfMIME != -1 {
			mime = mime[:endOfMIME]
		}
		mime = strings.TrimSpace(mime)

		mimes = append(mimes, mime)
	}

	return mimes
}
