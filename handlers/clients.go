package handlers

import (
	context "context"
	"io"

	"github.com/ONSdigital/dp-api-clients-go/v2/articles"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-renderer/model"
)

// To mock interfaces in this file
//go:generate mockgen -source=clients.go -destination=mock_clients.go -package=handlers github.com/ONSdigital/dp-frontend-articles-controller/handlers

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

// RenderClient is an interface with methods for rendering a template
type RenderClient interface {
	BuildPage(w io.Writer, pageModel interface{}, templateName string)
	NewBasePageModel() model.Page
}

// ZebedeeClient is an interface for zebedee client
type ZebedeeClient interface {
	GetBreadcrumb(ctx context.Context, userAccessToken, collectionID, lang, uri string) ([]zebedee.Breadcrumb, error)
	GetHomepageContent(ctx context.Context, userAccessToken, collectionID, lang, path string) (m zebedee.HomepageContent, err error)
}

// ArticlesApiClient is an interface for the Articles API client
type ArticlesApiClient interface {
	GetLegacyBulletin(ctx context.Context, userAccessToken, collectionID, lang, uri string) (*articles.Bulletin, error)
}
