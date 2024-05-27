package api

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Spanner controller allows methods to access data such as session cookies
// and environment variables such as API keys without exposiung them globally
type SpannerController struct {
	session *session.Store
}

func NewSpannerStorage(session *session.Store) *SpannerController {
	return &SpannerController{session: session}
}
