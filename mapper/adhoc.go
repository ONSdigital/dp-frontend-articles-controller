package mapper

import (
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-articles-controller/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

func CreateAdHocModel(basePage coreModel.Page, adhocJson model.AdHocJSON, bcs []zebedee.Breadcrumb, lang, requestProtocol string) model.AdHoc {
	adhoc := model.AdHoc{
		Page: basePage,
	}
	adhoc.Language = lang
	adhoc.BetaBannerEnabled = true
	adhoc.Type = adhocJson.Type
	adhoc.Metadata = coreModel.Metadata{
		Title:       adhocJson.Description.Title,
		Description: adhocJson.Description.MetaDescription,
		Keywords:    adhocJson.Description.Keywords,
	}
	adhoc.URI = adhocJson.Uri
	adhoc.ReleaseDate = adhocJson.Description.ReleaseDate
	adhoc.BodyMarkdown = adhocJson.Markdown
	adhoc.Reference = adhocJson.Description.Reference
	adhoc.Downloads = createDownloads(adhocJson.Downloads)
	return adhoc
}

func createDownloads(downloadsJSON []model.DownloadJSON) []model.Download {
	downloads := make([]model.Download, len(downloadsJSON))
	for index, download := range downloadsJSON {
		downloads[index] = model.Download(download)
	}
	return downloads
}
