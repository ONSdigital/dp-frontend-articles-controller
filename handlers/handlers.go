package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-articles-controller/assets"
	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/mapper"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-renderer/onshelper"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
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

	getContent := func(path string) (interface{}, error) {
		return zc.GetFigure(ctx, userAccessToken, collectionID, lang, path)
	}

	helper := onshelper.NewHelper(assets.Asset, assets.AssetNames, cfg.PatternLibraryAssetsPath, cfg.SiteDomain, getContent)

	basePage := rc.NewBasePageModel()
	model := mapper.CreateSixteensBulletinModel(basePage, *bulletin, breadcrumbs, lang)
	rc.BuildPageWithOptions(w, model, "sixteens-bulletin", helper.GetFuncMap())
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
	model := mapper.CreateBulletinModel(basePage, *bulletin, breadcrumbs, lang)
	rc.BuildPage(w, model, "bulletin")
}
