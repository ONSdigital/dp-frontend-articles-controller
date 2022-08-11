package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/mapper"
	"github.com/ONSdigital/dp-frontend-articles-controller/model"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
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
	//model := mapper.CreateBulletinModel(basePage, *bulletin, breadcrumbs, lang, requestProtocol)
	//rc.BuildPage(w, model, "bulletin")

	adhocJsonSource := []byte(`{
		"downloads": [
			{
				"title": "ANWAR",
				"file": "deccanwar2november2015_tcm77-430106.xls"
			},
			{
				"title": "Tables",
				"file": "decctablesnovember2015_tcm77-430108.pdf"
			}
		],
		"markdown": [
			"Export and import estimates of crude oil and other fuels including adjustments. "
		],
		"links": [],
		"type": "static_adhoc",
		"uri": "/economy/environmentalaccounts/adhocs/005204fuelandenergydataprovidedonamonthlybasistodecc",
		"description": {
			"title": "Fuel and energy data provided on a monthly basis to DECC",
			"keywords": [],
			"metaDescription": "Export and import estimates of crude oil and other fuels including adjustments. ",
			"releaseDate": "2016-01-15T08:03:43.411Z",
			"unit": "",
			"preUnit": "",
			"source": "",
			"reference": "005204"
		}
	}`)

	var adhocJson model.AdHocJSON
	jsonError := json.Unmarshal(adhocJsonSource, &adhocJson)
	if jsonError != nil {
		fmt.Println("error:", jsonError)
	}

	model := mapper.CreateAdHocModel(basePage, adhocJson, breadcrumbs, lang, requestProtocol)
	rc.BuildPage(w, model, "adhoc")
}
