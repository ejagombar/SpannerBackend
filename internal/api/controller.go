package api

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

type SpannerController struct {
	session *session.Store
}

func NewSpannerStorage(session *session.Store) *SpannerController {
	return &SpannerController{session: session}
}
