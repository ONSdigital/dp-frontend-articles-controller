package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/mapper"
	"github.com/ONSdigital/log.go/v2/log"
)

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

// Bulletin handles bulletin requests
func Bulletin(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bulletin(w, req, cfg)
	}
}

func bulletin(w http.ResponseWriter, req *http.Request, cfg config.Config) {
	ctx := req.Context()
	bulletin := mapper.Bulletin{Name: req.URL.EscapedPath()}
	model := mapper.Blank(ctx, bulletin, cfg)

	b, err := json.Marshal(model)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Error(ctx, "failed to write bytes for http response", err)
		setStatusCode(req, w, err)
		return
	}
	return
}
