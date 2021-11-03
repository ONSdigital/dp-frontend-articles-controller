package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/mapper"
	"github.com/ONSdigital/log.go/v2/log"
)

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
func Bulletin(cfg config.Config, rc RenderClient, zc ZebedeeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bulletin(w, req, rc, zc, cfg)
	}
}

func bulletin(w http.ResponseWriter, req *http.Request, rc RenderClient, zc ZebedeeClient, cfg config.Config) {
	ctx := req.Context()

	lang := "en"
	userAccessToken := ""

	bulletin, err := zc.GetBulletin(ctx, userAccessToken, lang, req.URL.EscapedPath())
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	breadcrumbs, err := zc.GetBreadcrumb(ctx, userAccessToken, "", lang, bulletin.URI)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	model := mapper.CreateBulletinModel(basePage, bulletin, breadcrumbs)
	rc.BuildPage(w, model, "bulletin")
}
