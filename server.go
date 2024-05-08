package daisy

import (
	"github.com/emersion/go-webdav/carddav"
	"net/http"
	"time"
)

type Conf struct {
	Listen string `json:"listen"`
}

func NewHttpServer(conf Conf, wd string) *http.Server {
	h := &carddav.Handler{
		Backend: &Backend{},
	}

	return &http.Server{
		Addr:              conf.Listen,
		Handler:           h,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    2500,
	}
}
