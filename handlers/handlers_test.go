package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/articles"
	"github.com/ONSdigital/dp-api-clients-go/v2/headers"
	zebedee "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	lang         = "en"
	accessToken  = "token"
	collectionID = "collection"
)

type testCliError struct{}

func (e *testCliError) Error() string { return "client error" }
func (e *testCliError) Code() int     { return http.StatusNotFound }

func TestUnitHandlers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := gomock.Any()

	Convey("test setStatusCode", t, func() {

		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "http://localhost:26500", nil)
			w := httptest.NewRecorder()
			err := &testCliError{}

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("test status code handles internal server error", func() {
			req := httptest.NewRequest("GET", "http://localhost:26500", nil)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test SixteensBulletin", t, func() {
		const requestUrlFormat = "http://localhost:26500/sixteens%s"
		url := "/a/bulletin/url"
		b := articles.Bulletin{
			URI:  "/the/bulletin/url",
			Type: "bulletin",
		}
		mockZebedeeClient := NewMockZebedeeClient(mockCtrl)
		mockRenderClient := NewMockRenderClient(mockCtrl)
		mockArticlesApiClient := NewMockArticlesApiClient(mockCtrl)
		mockConfig := config.Config{}

		router := mux.NewRouter()
		router.HandleFunc("/sixteens{uri:/.*}", SixteensBulletin(mockConfig, mockRenderClient, mockZebedeeClient, mockArticlesApiClient))

		w := httptest.NewRecorder()

		Convey("it returns 200 when rendered succesfully", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, accessToken, collectionID, lang, url).Return(&b, nil)
			mockZebedeeClient.EXPECT().GetBreadcrumb(ctx, accessToken, collectionID, lang, b.URI)
			mockRenderClient.EXPECT().NewBasePageModel()
			mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "sixteens-bulletin")

			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("it returns 200 when rendered succesfully without headers or cookies", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, "", "", lang, url).Return(&b, nil)
			mockZebedeeClient.EXPECT().GetBreadcrumb(ctx, "", "", lang, b.URI)
			mockRenderClient.EXPECT().NewBasePageModel()
			mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "sixteens-bulletin")

			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("it returns 500 when there is an error getting the bulletin from Zebedee", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, accessToken, collectionID, lang, url).Return(nil, errors.New(("error reading data")))

			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("it returns 500 when there is an error getting the breadcrumbs from Zebedee", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, accessToken, collectionID, lang, url).Return(&b, nil)
			mockZebedeeClient.EXPECT().GetBreadcrumb(ctx, accessToken, collectionID, lang, b.URI).Return([]zebedee.Breadcrumb{}, errors.New(("error reading breadcrumbs")))
			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test Bulletin", t, func() {
		requestUrlFormat := "http://localhost:26500%s"
		url := "/a/bulletin/url"
		b := articles.Bulletin{
			URI:  "/the/bulletin/url",
			Type: "bulletin",
		}
		mockZebedeeClient := NewMockZebedeeClient(mockCtrl)
		mockRenderClient := NewMockRenderClient(mockCtrl)
		mockArticlesApiClient := NewMockArticlesApiClient(mockCtrl)
		mockConfig := config.Config{}

		router := mux.NewRouter()
		router.HandleFunc(url, Bulletin(mockConfig, mockRenderClient, mockZebedeeClient, mockArticlesApiClient))

		w := httptest.NewRecorder()

		Convey("it returns 200 when rendered succesfully", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, accessToken, collectionID, lang, url).Return(&b, nil)
			mockZebedeeClient.EXPECT().GetBreadcrumb(ctx, accessToken, collectionID, lang, b.URI)
			mockZebedeeClient.EXPECT().GetHomepageContent(ctx, accessToken, collectionID, lang, "/")
			mockRenderClient.EXPECT().NewBasePageModel()
			mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "bulletin")

			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("it returns 200 when rendered succesfully without headers or cookies", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, "", "", lang, url).Return(&b, nil)
			mockZebedeeClient.EXPECT().GetBreadcrumb(ctx, "", "", lang, b.URI)
			mockZebedeeClient.EXPECT().GetHomepageContent(ctx, "", "", lang, "/")
			mockRenderClient.EXPECT().NewBasePageModel()
			mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "bulletin")

			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("it returns 500 when there is an error getting the bulletin from Zebedee", func() {
			mockZebedeeClient.EXPECT().GetHomepageContent(ctx, accessToken, collectionID, lang, "/")
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, accessToken, collectionID, lang, url).Return(nil, errors.New(("error reading data")))

			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("it returns 500 when there is an error getting the breadcrumbs from Zebedee", func() {
			mockZebedeeClient.EXPECT().GetHomepageContent(ctx, accessToken, collectionID, lang, "/")
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, accessToken, collectionID, lang, url).Return(&b, nil)
			mockZebedeeClient.EXPECT().GetBreadcrumb(ctx, accessToken, collectionID, lang, b.URI).Return([]zebedee.Breadcrumb{}, errors.New(("error reading breadcrumbs")))
			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test Bulletin '/data' endpoint", t, func() {
		requestUrlFormat := "http://localhost:26500%s"
		b := articles.Bulletin{
			URI:  "/a/bulletin/url",
			Type: "bulletin",
		}
		url := b.URI + "/data"

		mockArticlesApiClient := NewMockArticlesApiClient(mockCtrl)
		mockConfig := config.Config{}
		router := mux.NewRouter()
		router.HandleFunc(url, BulletinData(mockConfig, mockArticlesApiClient))
		w := httptest.NewRecorder()

		js, _ := json.Marshal(b)
		Convey("when the release is retrieved successfully", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, accessToken, collectionID, lang, b.URI).Return(&b, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			Convey("it returns 200 with the expected json payload ", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.Bytes(), ShouldResemble, js)
			})
			Convey("and the content type is 'application/json' ", func() {
				So(w.Header().Get("content-type"), ShouldEqual, "application/json")
			})
		})

		Convey("when the release is retrieved successfully without headers or cookies", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, "", "", lang, b.URI).Return(&b, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)

			router.ServeHTTP(w, req)

			Convey("it returns 200 with the expected json payload ", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.Bytes(), ShouldResemble, js)
			})
			Convey("and the content type is 'application/json' ", func() {
				So(w.Header().Get("content-type"), ShouldEqual, "application/json")
			})
		})

		Convey("it returns 500 when there is an error getting the release from the api", func() {
			mockArticlesApiClient.EXPECT().GetLegacyBulletin(ctx, accessToken, collectionID, lang, b.URI).Return(nil, errors.New("error reading data"))
			req := httptest.NewRequest("GET", fmt.Sprintf(requestUrlFormat, url), nil)
			setRequestHeaders(req)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func setRequestHeaders(req *http.Request) {
	headers.SetAuthToken(req, accessToken)
	headers.SetCollectionID(req, collectionID)
}
