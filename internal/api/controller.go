package api

//
import (
	"github.com/ejagombar/SpannerBackend/config"
)

type SpannerController struct {
	storage *SpannerStorage
	config  *config.EnvVars
}

func NewSpannerController(storage *SpannerStorage, config *config.EnvVars) *SpannerController {
	return &SpannerController{storage: storage, config: config}
}
