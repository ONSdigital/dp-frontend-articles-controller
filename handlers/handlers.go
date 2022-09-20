package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/mapper"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

const homepagePath = "/"

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
func SixteensBulletin(cfg config.Config, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		sixteensBulletin(w, r, accessToken, collectionID, lang, rc, zc, ac, cfg)
	})
}

func sixteensBulletin(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient, cfg config.Config) {
	ctx := req.Context()
	muxVars := mux.Vars(req)
	uri := muxVars["uri"]
	bulletin, err := ac.GetLegacyBulletin(ctx, userAccessToken, collectionID, lang, uri)

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
	model := mapper.CreateSixteensBulletinModel(basePage, *bulletin, breadcrumbs, lang)
	rc.BuildPage(w, model, "sixteens-bulletin")
}

// Bulletin handles bulletin requests
func Bulletin(cfg config.Config, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		bulletin(w, r, accessToken, collectionID, lang, rc, zc, ac, cfg)
	})
}

func bulletin(w http.ResponseWriter, req *http.Request, userAccessToken, collectionID, lang string, rc RenderClient, zc ZebedeeClient, ac ArticlesApiClient, cfg config.Config) {
	ctx := req.Context()

	homepageContent, err := zc.GetHomepageContent(ctx, userAccessToken, collectionID, lang, homepagePath)
	if err != nil {
		log.Warn(ctx, "unable to get homepage content", log.FormatErrors([]error{err}), log.Data{"homepage_content": err})
	}

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
	requestProtocol := "http"
	if req.TLS != nil {
		requestProtocol = "https"
	}
	model := mapper.CreateBulletinModel(basePage, *bulletin, breadcrumbs, lang, requestProtocol, homepageContent.ServiceMessage, homepageContent.EmergencyBanner)
	rc.BuildPage(w, model, "bulletin")
}
