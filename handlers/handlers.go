package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/mapper"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
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
func Bulletin(cfg config.Config, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		bulletin(w, r, accessToken, collectionID, lang, rc, zc, ac, cfg)
	})
}

func bulletin(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient, cfg config.Config) {
	ctx := req.Context()

	bulletin, err := ac.GetLegacyBulletin(ctx, userAccessToken, collectionID, lang, req.URL.EscapedPath())

	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	breadcrumbs, err := zc.GetBreadcrumb(ctx, userAccessToken, collectionID, lang, bulletin.URI)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	basePage := rc.NewBasePageModel()
	model := mapper.CreateBulletinModel(basePage, *bulletin, breadcrumbs)
	rc.BuildPage(w, model, "bulletin")
}
