package api

//
import (
	"github.com/ejagombar/SpannerBackend/config"
	"github.com/zmb3/spotify/v2"
)

type SpannerController struct {
	storage *SpannerStorage
	config  *config.EnvVars
	client  *spotify.Client
}

func NewSpannerController(storage *SpannerStorage, config *config.EnvVars) *SpannerController {
	return &SpannerController{storage: storage, config: config}
}
