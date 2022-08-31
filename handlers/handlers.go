package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-frontend-articles-controller/assets"
	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/mapper"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	"github.com/ONSdigital/dp-renderer/sixteenstagresolver"
	"github.com/ONSdigital/dp-renderer/tagresolver"
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

	resourceReader := sixteenstagresolver.ResourceReader{
		GetFigure: func(path string) (coreModel.Figure, error) {
			fmt.Printf("GetFigure path %s\n", path)
			figure, err := zc.GetFigure(ctx, userAccessToken, collectionID, lang, path)
			if err != nil {
				return coreModel.Figure{}, err
			}
			fmt.Printf("GetFigure unmapped %#v\n\n", figure)
			return mapper.MapFigure(figure), nil
		},
		GetResourceBody: func(path string) ([]byte, error) {
			return zc.GetResourceBody(ctx, userAccessToken, collectionID, lang, path)
		},
		GetTable: func(html []byte) (string, error) {
			//TODO create table renderer client
			//Call table renderer with html
			tableRendererUrl := "http://localhost:23300"
			req, err := http.NewRequest("POST", tableRendererUrl+"/render/html", bytes.NewBuffer(html))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			return string(body), nil
		},
		GetFileSize: func(path string) (int, error) {
			size, err := zc.GetFileSize(ctx, userAccessToken, collectionID, lang, path)
			if err != nil {
				return 0, err
			}
			return size.Size, nil
		},
	}

	resolverCfg := sixteenstagresolver.TagResolverRenderConfig{
		Asset:                    assets.Asset,
		AssetNames:               assets.AssetNames,
		PatternLibraryAssetsPath: cfg.PatternLibraryAssetsPath,
		SiteDomain:               cfg.SiteDomain,
	}

	helper := sixteenstagresolver.NewTagResolverHelper(uri, resourceReader, resolverCfg)

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

	resourceReader := tagresolver.ResourceReader{
		GetFigure: func(path string) (coreModel.Figure, error) {
			fmt.Printf("GetFigure path %s\n", path)
			figure, err := zc.GetFigure(ctx, userAccessToken, collectionID, lang, path)
			if err != nil {
				return coreModel.Figure{}, err
			}
			fmt.Printf("GetFigure unmapped %#v\n\n", figure)
			return mapper.MapFigure(figure), nil
		},
		GetResourceBody: func(path string) ([]byte, error) {
			return zc.GetResourceBody(ctx, userAccessToken, collectionID, lang, path)
		},
		GetTable: func(html []byte) (string, error) {
			//TODO create table renderer client
			//Call table renderer with html
			tableRendererUrl := "http://localhost:23300"
			req, err := http.NewRequest("POST", tableRendererUrl+"/render/html", bytes.NewBuffer(html))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			return string(body), nil
		},
		GetFileSize: func(path string) (int, error) {
			size, err := zc.GetFileSize(ctx, userAccessToken, collectionID, lang, path)
			if err != nil {
				return 0, err
			}
			return size.Size, nil
		},
	}

	resolverCfg := tagresolver.TagResolverRenderConfig{
		Asset:                    assets.Asset,
		AssetNames:               assets.AssetNames,
		PatternLibraryAssetsPath: cfg.PatternLibraryAssetsPath,
		SiteDomain:               cfg.SiteDomain,
	}

	helper := tagresolver.NewTagResolverHelper(uri, resourceReader, resolverCfg)

	basePage := rc.NewBasePageModel()
	requestProtocol := "http"
	if req.TLS != nil {
		requestProtocol = "https"
	}
	model := mapper.CreateBulletinModel(basePage, *bulletin, breadcrumbs, lang, requestProtocol)
	rc.BuildPageWithOptions(w, model, "bulletin", helper.GetFuncMap())
}
