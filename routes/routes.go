package routes

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/articles"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-articles-controller/config"
	"github.com/ONSdigital/dp-frontend-articles-controller/handlers"
	render "github.com/ONSdigital/dp-renderer"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	HealthCheckHandler func(w http.ResponseWriter, req *http.Request)
	Zebedee            *zebedee.Client
	Render             *render.Render
	ArticlesAPI        *articles.Client
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, c Clients) {
	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(c.HealthCheckHandler)
	r.StrictSlash(true).Path("/sixteens/{uri:.*}").Methods("GET").HandlerFunc(handlers.SixteensBulletin(*cfg, c.Render, c.Zebedee, c.ArticlesAPI))
	r.StrictSlash(true).Path("/{uri:.*}").Methods("GET").HandlerFunc(handlers.Bulletin(*cfg, c.Render, c.Zebedee, c.ArticlesAPI))
}
