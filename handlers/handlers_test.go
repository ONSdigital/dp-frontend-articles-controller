package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	zebedee "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
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

	Convey("test GetBulletin", t, func() {
		url := "/a/bulletin/url"
		mockZebedeeClient := NewMockZebedeeClient(mockCtrl)
		mockRenderClient := NewMockRenderClient(mockCtrl)
		mockConfig := config.Config{}

		router := mux.NewRouter()
		router.HandleFunc(url, Bulletin(mockConfig, mockRenderClient, mockZebedeeClient))

		w := httptest.NewRecorder()

		Convey("it returns 200 when rendered succesfully", func() {
			mockZebedeeClient.EXPECT().GetBulletin(ctx, "", "en", url)
			mockRenderClient.EXPECT().NewBasePageModel()
			mockRenderClient.EXPECT().BuildPage(w, gomock.Any(), "bulletin") // TODO verify interaction

			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:26500%s", url), nil)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("it returns 500 when there is an error getting the bulletin from Zebedee", func() {
			mockZebedeeClient.EXPECT().GetBulletin(ctx, "", "en", url).Return(zebedee.Bulletin{}, errors.New(("error reading data")))

			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:26500%s", url), nil)

			router.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}
