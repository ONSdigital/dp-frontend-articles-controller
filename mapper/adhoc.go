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
	adhoc.TableOfContents = createAdHocTableOfContents()
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

func createAdHocTableOfContents() coreModel.TableOfContents {
	toc := coreModel.TableOfContents{
		Id: "toc",
		AriaLabel: coreModel.Localisation{
			LocaleKey: "TableOfContents",
			Plural:    1,
		},
		Title: coreModel.Localisation{
			LocaleKey: "PageContentsTitle",
			Plural:    1,
		},
	}

	sections := make(map[string]coreModel.ContentSection)
	sections["section-summary"] = coreModel.ContentSection{
		Current: false,
		Title: coreModel.Localisation{
			LocaleKey: "SectionTitleSummaryOfRequest",
			Plural:    1,
		},
	}
	sections["section-downloads"] = coreModel.ContentSection{
		Current: false,
		Title: coreModel.Localisation{
			LocaleKey: "SectionTitleDownloadAssociatedWithRequest",
			Plural:    1,
		},
	}
	sections["section-further-reading"] = coreModel.ContentSection{
		Current: false,
		Title: coreModel.Localisation{
			LocaleKey: "SectionTitleFurtherReading",
			Plural:    1,
		},
	}

	displayOrder := []string{
		"section-summary",
		"section-downloads",
		"section-further-reading",
	}

	toc.Sections = sections
	toc.DisplayOrder = displayOrder

	return toc
}
