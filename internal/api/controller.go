package api

//
import (
// "github.com/gofiber/fiber/v2/middleware/session"
)

type SpannerController struct {
	storage *SpannerStorage
}

func NewSpannerController(storage *SpannerStorage) *SpannerController {
	return &SpannerController{storage: storage}
}
