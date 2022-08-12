package handlers

import (
	"bytes"
	context "context"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
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

	// replaceCharts(ctx, userAccessToken, collectionID, lang, zc, bulletin, cfg)

	// STEVE's

	getContent := func(path string) (interface{}, error) {
		return zc.GetFigure(ctx, userAccessToken, collectionID, lang, path)
	}

	// TODO: resolve development flag here
	helper := onshelper.NewHelper(assets.Asset, assets.AssetNames, cfg.PatternLibraryAssetsPath, cfg.SiteDomain, getContent)

	//

	basePage := rc.NewBasePageModel()
	model := mapper.CreateSixteensBulletinModel(basePage, *bulletin, breadcrumbs, lang)
	rc.BuildPageWithOptions(w, model, "sixteens-bulletin", helper.GetFuncMap())
}

type FigureReplacement struct {
	regex    *regexp.Regexp
	template string
}

func replaceCharts(ctx context.Context, userAccessToken, collectionID, lang string, zc ZebedeeClient, bulletin *zebedee.Bulletin, cfg config.Config) {
	figureReplacements := []FigureReplacement{
		{
			regex:    regexp.MustCompile(`<ons-chart path=\"(.*)\" />`),
			template: "partials/chart",
		},
	}

	// rc := render.New(client.NewUnrolledAdapterForPartials(assets.Asset, assets.AssetNames, true), cfg.PatternLibraryAssetsPath, cfg.SiteDomain)

	// Concurrently resolve figure data coming from zebedee
	var wg sync.WaitGroup
	// We use this buffered channel to limit the number of concurrent calls we make to zebedee
	sem := make(chan int, 10)

	for i, s := range bulletin.Sections {
		for _, rep := range figureReplacements {
			matches := rep.regex.FindAllStringSubmatch(s.Markdown, -1)
			for _, m := range matches {
				sem <- 1
				wg.Add(1)
				go func(i int, figureTag, figureUrl string) {
					defer func() {
						<-sem
						wg.Done()
					}()

					buf := new(bytes.Buffer)
					// model, _ := zc.GetFigure(ctx, userAccessToken, collectionID, lang, figureUrl)
					// rc.BuildPage(buf, model, rep.template)
					bulletin.Sections[i].Markdown = strings.Replace(bulletin.Sections[i].Markdown, figureTag, buf.String(), 1)
				}(i, m[0], m[1])
			}
		}
	}

	wg.Wait()
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
