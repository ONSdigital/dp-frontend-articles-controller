package mapper

import (
	"html/template"
	"net/url"
	"sort"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/articles"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-articles-controller/mocks"
	"github.com/ONSdigital/dp-renderer/helper"
	"github.com/ONSdigital/dp-renderer/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)

	Convey("Given a bulletin, basePage and breadcrumbs", t, func() {
		basePage := coreModel.NewPage("path/to/assets", "site-domain")

		bulletin := articles.Bulletin{
			Type: "bulletin",
			Description: zebedee.Description{
				Title:             "Title",
				Edition:           "2021",
				Keywords:          []string{"energy", "waste"},
				MetaDescription:   "description",
				Summary:           "summary",
				NationalStatistic: true,
				Contact: zebedee.Contact{
					Name:      "contact name",
					Email:     "contact@ons.gov.uk",
					Telephone: "+44",
				},
				ReleaseDate:   "2015-07-08T23:00:00.000Z",
				NextRelease:   "",
				LatestRelease: true,
				DatasetID:     "22",
			},
			LatestReleaseURI: "uri/2022",
			Sections: []zebedee.Section{
				{
					Title:    "section1",
					Markdown: "markdown1",
				},
				{
					Title:    "section2",
					Markdown: "markdown2",
				},
			},
			Accordion: []zebedee.Section{
				{
					Title:    "acco1",
					Markdown: "accordion markdown1",
				},
				{
					Title:    "acco2",
					Markdown: "accordion markdown2",
				},
			},
			RelatedBulletins: []zebedee.Link{
				{
					Title: "related bulletin 1",
					URI:   "bulletin link 1",
				},
				{
					Title: "related bulletin 2",
					URI:   "bulletin link 2",
				},
			},
			RelatedData: []zebedee.Link{
				{
					Title: "related data 1",
					URI:   "data link 1",
				},
				{
					Title: "related 2",
					URI:   "data link 2",
				},
			},
			Links: []zebedee.Link{
				{
					Title: "link 1",
					URI:   "http 1",
				},
				{
					Title: "link 2",
					URI:   "http 2",
				},
			},
			Charts: []zebedee.Figure{
				{
					Title:    "chart 1",
					Filename: "chart1.png",
					Version:  "1",
					URI:      "chart 1 uri",
				},
				{
					Title:    "chart 2",
					Filename: "chart2.png",
					Version:  "2",
					URI:      "chart 2 uri",
				},
			},
			Tables: []zebedee.Figure{
				{
					Title:    "table 1",
					Filename: "table1.png",
					Version:  "1",
					URI:      "table 1 uri",
				},
				{
					Title:    "table 2",
					Filename: "table2.png",
					Version:  "2",
					URI:      "table 2 uri",
				},
			},
			Images: []zebedee.Figure{
				{
					Title:    "image 1",
					Filename: "img1.png",
					Version:  "1",
					URI:      "img 1 uri",
				},
				{
					Title:    "image 2",
					Filename: "img2.png",
					Version:  "2",
					URI:      "img 2 uri",
				},
			},
			Equations: []zebedee.Figure{
				{
					Title:    "eq 1",
					Filename: "eq.png",
					Version:  "1",
					URI:      "eq 1 uri",
				},
				{
					Title:    "eq 2",
					Filename: "eq2.png",
					Version:  "2",
					URI:      "eq 2 uri",
				},
			},
			Versions: []zebedee.Version{
				{
					ReleaseDate: "2021-07-19T21:42:11.302Z",
					Notice:      "notice2",
					URI:         "/n2",
				},
				{
					ReleaseDate: "2021-10-19T10:43:34.507Z",
					Notice:      "notice1",
					URI:         "/n1",
				},
			},
		}

		breadcrumbs := []zebedee.Breadcrumb{
			{
				Description: zebedee.NodeDescription{
					Title: "Home",
				},
				URI: "/",
			},
			{
				Description: zebedee.NodeDescription{
					Title: "Economy",
				},
				URI: "/economy",
			},
			{
				Description: zebedee.NodeDescription{
					Title: "Product title",
				},
				URI: "/economy/product",
			},
		}

		serviceMessage := "Service Message"
		emergencyBannerTitle := "Emergency Title"
		emergencyBannerType := "notable-death"
		emergencyBannerDescription := "Emergency Description"
		emergencyBannerUri := "https://example.com/emergency"
		emergencyBannerLinkText := "Attention, this is an emergency. There's an emergency going on."
		bannerData := zebedee.EmergencyBanner{
			Title:       emergencyBannerTitle,
			Type:        emergencyBannerType,
			Description: emergencyBannerDescription,
			URI:         emergencyBannerUri,
			LinkText:    emergencyBannerLinkText,
		}

		Convey("When the bulletin belongs to a census survey", func() {
			bulletin.Description.Survey = "census"

			Convey("And the bulletin URI is not a previous version", func() {
				bulletin.URI = "the/bulletin/uri/path/version"

				Convey("CreateBulletinModel maps correctly", func() {
					requestProtocol := "https"
					model := CreateBulletinModel(basePage, bulletin, breadcrumbs, "cy", requestProtocol, serviceMessage, bannerData)

					So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
					So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
					So(model.BetaBannerEnabled, ShouldBeTrue)
					So(model.BetaBannerEnabled, ShouldBeTrue)
					So(model.ServiceMessage, ShouldEqual, serviceMessage)
					So(model.EmergencyBanner.Title, ShouldEqual, emergencyBannerTitle)
					So(model.EmergencyBanner.Type, ShouldEqual, emergencyBannerType)
					So(model.EmergencyBanner.Description, ShouldEqual, emergencyBannerDescription)
					So(model.EmergencyBanner.URI, ShouldEqual, emergencyBannerUri)
					So(model.EmergencyBanner.LinkText, ShouldEqual, emergencyBannerLinkText)
					So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
					So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
					So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
					So(model.Summary, ShouldEqual, bulletin.Description.Summary)
					So(model.Type, ShouldEqual, bulletin.Type)
					So(model.URI, ShouldEqual, bulletin.URI)
					So(model.ParentPath, ShouldEqual, "the/bulletin/uri/path")
					So(model.CorrectedPath, ShouldBeEmpty)
					So(model.Edition, ShouldEqual, bulletin.Description.Edition)
					So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
					So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
					So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
					So(model.LatestRelease, ShouldBeTrue)
					So(model.LatestReleaseUri, ShouldEqual, bulletin.LatestReleaseURI)
					So(model.DatasetId, ShouldEqual, bulletin.Description.DatasetID)
					So(model.Census2021, ShouldEqual, true)
					So(model.AboutTheData, ShouldEqual, true)
					So(model.PreGTMJavaScript, ShouldNotBeEmpty)
					So(len(model.PreGTMJavaScript), ShouldEqual, 1)
					expectedPreGTMJs := createPreGTMJs(bulletin.Description.Title, bulletin.Description.ReleaseDate, bulletin.URI, "census")
					So(model.PreGTMJavaScript[0], ShouldResemble, expectedPreGTMJs)
					So(len(model.Sections), ShouldEqual, len(bulletin.Sections))
					assertSections(model.Sections, bulletin.Sections)
					assertSections(model.Accordion, bulletin.Accordion)
					assertContentsView(model.ContentsView, bulletin.Sections, bulletin.Accordion, model.AboutTheData)
					assertLinks(model.RelatedBulletins, bulletin.RelatedBulletins)
					assertLinks(model.RelatedData, bulletin.RelatedData)
					assertLinks(model.Links, bulletin.Links)
					assertFigures(model.Charts, bulletin.Charts)
					assertFigures(model.Tables, bulletin.Tables)
					assertFigures(model.Images, bulletin.Images)
					assertFigures(model.Equations, bulletin.Equations)
					assertShareLinks(model.ShareLinks, bulletin.URI, requestProtocol)
					So(len(model.Versions), ShouldEqual, len(bulletin.Versions))
					// Versions should be sorted by date
					sort.Slice(bulletin.Versions, func(i, j int) bool { return bulletin.Versions[i].ReleaseDate > bulletin.Versions[j].ReleaseDate })
					for i, v := range bulletin.Versions {
						f := model.Versions[i]
						So(f.Date, ShouldEqual, v.ReleaseDate)
						So(f.Markdown, ShouldEqual, v.Notice)
						So(f.URI, ShouldEqual, v.URI)
					}
					So(len(model.Alerts), ShouldEqual, len(bulletin.Alerts))
					// Alerts should be sorted by date
					sort.Slice(bulletin.Alerts, func(i, j int) bool { return bulletin.Alerts[i].Date > bulletin.Alerts[j].Date })
					for i, a := range bulletin.Alerts {
						f := model.Alerts[i]
						So(f.Date, ShouldEqual, a.Date)
						So(f.Markdown, ShouldEqual, a.Markdown)
					}
					So(len(model.Breadcrumb), ShouldEqual, len(breadcrumbs))
					for i, b := range breadcrumbs {
						found := model.Breadcrumb[i]
						if i == 0 && b.Description.Title == "Home" {
							So(found.Title, ShouldEqual, "Hafan")
						} else {
							So(found.Title, ShouldEqual, b.Description.Title)
						}
						So(found.URI, ShouldEqual, b.URI)
					}
				})

				Convey("CreateSixteensBulletinModel maps correctly", func() {
					model := CreateSixteensBulletinModel(basePage, bulletin, breadcrumbs, "cy")

					So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
					So(model.FeatureFlags.SixteensVersion, ShouldEqual, "67f6982")
					So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
					So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
					So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
					So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
					So(model.Summary, ShouldEqual, bulletin.Description.Summary)
					So(model.Type, ShouldEqual, bulletin.Type)
					So(model.URI, ShouldEqual, bulletin.URI)
					So(model.ParentPath, ShouldEqual, "the/bulletin/uri/path")
					So(model.CorrectedPath, ShouldBeEmpty)
					So(model.Edition, ShouldEqual, bulletin.Description.Edition)
					So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
					So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
					So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
					So(model.LatestRelease, ShouldBeTrue)
					So(model.LatestReleaseUri, ShouldEqual, bulletin.LatestReleaseURI)
					So(model.DatasetId, ShouldEqual, bulletin.Description.DatasetID)
					So(len(model.Sections), ShouldEqual, len(bulletin.Sections))
					assertSections(model.Sections, bulletin.Sections)
					assertSections(model.Accordion, bulletin.Accordion)
					assertLinks(model.RelatedBulletins, bulletin.RelatedBulletins)
					assertLinks(model.RelatedData, bulletin.RelatedData)
					assertLinks(model.Links, bulletin.Links)
					assertFigures(model.Charts, bulletin.Charts)
					assertFigures(model.Tables, bulletin.Tables)
					assertFigures(model.Images, bulletin.Images)
					assertFigures(model.Equations, bulletin.Equations)
					So(len(model.Versions), ShouldEqual, len(bulletin.Versions))
					// Versions should be sorted by date
					sort.Slice(bulletin.Versions, func(i, j int) bool { return bulletin.Versions[i].ReleaseDate > bulletin.Versions[j].ReleaseDate })
					for i, v := range bulletin.Versions {
						f := model.Versions[i]
						So(f.Date, ShouldEqual, v.ReleaseDate)
						So(f.Markdown, ShouldEqual, v.Notice)
						So(f.URI, ShouldEqual, v.URI)
					}
					So(len(model.Alerts), ShouldEqual, len(bulletin.Alerts))
					// Alerts should be sorted by date
					sort.Slice(bulletin.Alerts, func(i, j int) bool { return bulletin.Alerts[i].Date > bulletin.Alerts[j].Date })
					for i, a := range bulletin.Alerts {
						f := model.Alerts[i]
						So(f.Date, ShouldEqual, a.Date)
						So(f.Markdown, ShouldEqual, a.Markdown)
					}
					So(len(model.Breadcrumb), ShouldEqual, len(breadcrumbs))
					for i, b := range breadcrumbs {
						found := model.Breadcrumb[i]
						So(found.Title, ShouldEqual, b.Description.Title)
						So(found.URI, ShouldEqual, b.URI)
					}
				})
			})
			Convey("And the bulletin URI is a previous version", func() {
				bulletin.URI = "the/bulletin/uri/path/previous/version"
				Convey("CreateBulletinModel maps correctly", func() {
					requestProtocol := "https"
					model := CreateBulletinModel(basePage, bulletin, breadcrumbs, "cy", requestProtocol, serviceMessage, bannerData)

					So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
					So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
					So(model.BetaBannerEnabled, ShouldBeTrue)
					So(model.ServiceMessage, ShouldEqual, serviceMessage)
					So(model.EmergencyBanner.Title, ShouldEqual, emergencyBannerTitle)
					So(model.EmergencyBanner.Type, ShouldEqual, emergencyBannerType)
					So(model.EmergencyBanner.Description, ShouldEqual, emergencyBannerDescription)
					So(model.EmergencyBanner.URI, ShouldEqual, emergencyBannerUri)
					So(model.EmergencyBanner.LinkText, ShouldEqual, emergencyBannerLinkText)
					So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
					So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
					So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
					So(model.Type, ShouldEqual, bulletin.Type)
					So(model.URI, ShouldEqual, bulletin.URI)
					So(model.ParentPath, ShouldEqual, "the/bulletin/uri/path/previous")
					So(model.CorrectedPath, ShouldEqual, "the/bulletin/uri/path")
					So(model.Edition, ShouldEqual, bulletin.Description.Edition)
					So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
					So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
					So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
					So(model.LatestRelease, ShouldBeTrue)
					So(model.LatestReleaseUri, ShouldEqual, bulletin.LatestReleaseURI)
					So(model.DatasetId, ShouldEqual, bulletin.Description.DatasetID)
					So(model.Census2021, ShouldEqual, true)
					So(model.AboutTheData, ShouldEqual, true)
					So(model.PreGTMJavaScript, ShouldNotBeEmpty)
					So(len(model.PreGTMJavaScript), ShouldEqual, 1)
					expectedPreGTMJs := createPreGTMJs(bulletin.Description.Title, bulletin.Description.ReleaseDate, bulletin.URI, "census")
					So(model.PreGTMJavaScript[0], ShouldResemble, expectedPreGTMJs)
					So(len(model.Sections), ShouldEqual, len(bulletin.Sections))
					assertSections(model.Sections, bulletin.Sections)
					assertSections(model.Accordion, bulletin.Accordion)
					assertContentsView(model.ContentsView, bulletin.Sections, bulletin.Accordion, model.AboutTheData)
					assertLinks(model.RelatedBulletins, bulletin.RelatedBulletins)
					assertLinks(model.RelatedData, bulletin.RelatedData)
					assertLinks(model.Links, bulletin.Links)
					assertFigures(model.Charts, bulletin.Charts)
					assertFigures(model.Tables, bulletin.Tables)
					assertFigures(model.Images, bulletin.Images)
					assertFigures(model.Equations, bulletin.Equations)
					assertShareLinks(model.ShareLinks, bulletin.URI, requestProtocol)
					So(len(model.Versions), ShouldEqual, len(bulletin.Versions))
					// Versions should be sorted by date
					sort.Slice(bulletin.Versions, func(i, j int) bool { return bulletin.Versions[i].ReleaseDate > bulletin.Versions[j].ReleaseDate })
					for i, v := range bulletin.Versions {
						f := model.Versions[i]
						So(f.Date, ShouldEqual, v.ReleaseDate)
						So(f.Markdown, ShouldEqual, v.Notice)
						So(f.URI, ShouldEqual, v.URI)
					}
					So(len(model.Alerts), ShouldEqual, len(bulletin.Alerts))
					// Alerts should be sorted by date
					sort.Slice(bulletin.Alerts, func(i, j int) bool { return bulletin.Alerts[i].Date > bulletin.Alerts[j].Date })
					for i, a := range bulletin.Alerts {
						f := model.Alerts[i]
						So(f.Date, ShouldEqual, a.Date)
						So(f.Markdown, ShouldEqual, a.Markdown)
					}
					So(len(model.Breadcrumb), ShouldEqual, len(breadcrumbs))
					for i, b := range breadcrumbs {
						found := model.Breadcrumb[i]
						if i == 0 && b.Description.Title == "Home" {
							So(found.Title, ShouldEqual, "Hafan")
						} else {
							So(found.Title, ShouldEqual, b.Description.Title)
						}
						So(found.URI, ShouldEqual, b.URI)
					}
				})

				Convey("CreateSixteensBulletinModel maps correctly", func() {
					model := CreateSixteensBulletinModel(basePage, bulletin, breadcrumbs, "cy")

					So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
					So(model.FeatureFlags.SixteensVersion, ShouldEqual, "67f6982")
					So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
					So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
					So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
					So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
					So(model.Type, ShouldEqual, bulletin.Type)
					So(model.URI, ShouldEqual, bulletin.URI)
					So(model.ParentPath, ShouldEqual, "the/bulletin/uri/path/previous")
					So(model.CorrectedPath, ShouldEqual, "the/bulletin/uri/path")
					So(model.Edition, ShouldEqual, bulletin.Description.Edition)
					So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
					So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
					So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
					So(model.LatestRelease, ShouldBeTrue)
					So(model.LatestReleaseUri, ShouldEqual, bulletin.LatestReleaseURI)
					So(model.DatasetId, ShouldEqual, bulletin.Description.DatasetID)
					So(len(model.Sections), ShouldEqual, len(bulletin.Sections))
					assertSections(model.Sections, bulletin.Sections)
					assertSections(model.Accordion, bulletin.Accordion)
					assertLinks(model.RelatedBulletins, bulletin.RelatedBulletins)
					assertLinks(model.RelatedData, bulletin.RelatedData)
					assertLinks(model.Links, bulletin.Links)
					assertFigures(model.Charts, bulletin.Charts)
					assertFigures(model.Tables, bulletin.Tables)
					assertFigures(model.Images, bulletin.Images)
					assertFigures(model.Equations, bulletin.Equations)
					So(len(model.Versions), ShouldEqual, len(bulletin.Versions))
					// Versions should be sorted by date
					sort.Slice(bulletin.Versions, func(i, j int) bool { return bulletin.Versions[i].ReleaseDate > bulletin.Versions[j].ReleaseDate })
					for i, v := range bulletin.Versions {
						f := model.Versions[i]
						So(f.Date, ShouldEqual, v.ReleaseDate)
						So(f.Markdown, ShouldEqual, v.Notice)
						So(f.URI, ShouldEqual, v.URI)
					}
					So(len(model.Alerts), ShouldEqual, len(bulletin.Alerts))
					// Alerts should be sorted by date
					sort.Slice(bulletin.Alerts, func(i, j int) bool { return bulletin.Alerts[i].Date > bulletin.Alerts[j].Date })
					for i, a := range bulletin.Alerts {
						f := model.Alerts[i]
						So(f.Date, ShouldEqual, a.Date)
						So(f.Markdown, ShouldEqual, a.Markdown)
					}
					So(len(model.Breadcrumb), ShouldEqual, len(breadcrumbs))
					for i, b := range breadcrumbs {
						found := model.Breadcrumb[i]
						So(found.Title, ShouldEqual, b.Description.Title)
						So(found.URI, ShouldEqual, b.URI)
					}
				})
			})
		})

		Convey("When the bulletin does not belong to a census survey", func() {
			bulletin.Description.Survey = "other"

			Convey("And the bulletin URI is not a previous version", func() {
				bulletin.URI = "the/bulletin/uri/path/version"

				Convey("CreateBulletinModel maps correctly", func() {
					requestProtocol := "https"
					model := CreateBulletinModel(basePage, bulletin, breadcrumbs, "cy", requestProtocol, serviceMessage, bannerData)

					So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
					So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
					So(model.BetaBannerEnabled, ShouldBeTrue)
					So(model.BetaBannerEnabled, ShouldBeTrue)
					So(model.ServiceMessage, ShouldEqual, serviceMessage)
					So(model.EmergencyBanner.Title, ShouldEqual, emergencyBannerTitle)
					So(model.EmergencyBanner.Type, ShouldEqual, emergencyBannerType)
					So(model.EmergencyBanner.Description, ShouldEqual, emergencyBannerDescription)
					So(model.EmergencyBanner.URI, ShouldEqual, emergencyBannerUri)
					So(model.EmergencyBanner.LinkText, ShouldEqual, emergencyBannerLinkText)
					So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
					So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
					So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
					So(model.Summary, ShouldEqual, bulletin.Description.Summary)
					So(model.Type, ShouldEqual, bulletin.Type)
					So(model.URI, ShouldEqual, bulletin.URI)
					So(model.ParentPath, ShouldEqual, "the/bulletin/uri/path")
					So(model.CorrectedPath, ShouldBeEmpty)
					So(model.Edition, ShouldEqual, bulletin.Description.Edition)
					So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
					So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
					So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
					So(model.LatestRelease, ShouldBeTrue)
					So(model.LatestReleaseUri, ShouldEqual, bulletin.LatestReleaseURI)
					So(model.DatasetId, ShouldEqual, bulletin.Description.DatasetID)
					So(model.Census2021, ShouldEqual, false)
					So(model.AboutTheData, ShouldEqual, false)
					So(model.PreGTMJavaScript, ShouldNotBeEmpty)
					So(len(model.PreGTMJavaScript), ShouldEqual, 1)
					expectedPreGTMJs := createPreGTMJs(bulletin.Description.Title, bulletin.Description.ReleaseDate, bulletin.URI, "")
					So(model.PreGTMJavaScript[0], ShouldResemble, expectedPreGTMJs)
					So(len(model.Sections), ShouldEqual, len(bulletin.Sections))
					assertSections(model.Sections, bulletin.Sections)
					assertSections(model.Accordion, bulletin.Accordion)
					assertContentsView(model.ContentsView, bulletin.Sections, bulletin.Accordion, model.AboutTheData)
					assertLinks(model.RelatedBulletins, bulletin.RelatedBulletins)
					assertLinks(model.RelatedData, bulletin.RelatedData)
					assertLinks(model.Links, bulletin.Links)
					assertFigures(model.Charts, bulletin.Charts)
					assertFigures(model.Tables, bulletin.Tables)
					assertFigures(model.Images, bulletin.Images)
					assertFigures(model.Equations, bulletin.Equations)
					assertShareLinks(model.ShareLinks, bulletin.URI, requestProtocol)
					So(len(model.Versions), ShouldEqual, len(bulletin.Versions))
					// Versions should be sorted by date
					sort.Slice(bulletin.Versions, func(i, j int) bool { return bulletin.Versions[i].ReleaseDate > bulletin.Versions[j].ReleaseDate })
					for i, v := range bulletin.Versions {
						f := model.Versions[i]
						So(f.Date, ShouldEqual, v.ReleaseDate)
						So(f.Markdown, ShouldEqual, v.Notice)
						So(f.URI, ShouldEqual, v.URI)
					}
					So(len(model.Alerts), ShouldEqual, len(bulletin.Alerts))
					// Alerts should be sorted by date
					sort.Slice(bulletin.Alerts, func(i, j int) bool { return bulletin.Alerts[i].Date > bulletin.Alerts[j].Date })
					for i, a := range bulletin.Alerts {
						f := model.Alerts[i]
						So(f.Date, ShouldEqual, a.Date)
						So(f.Markdown, ShouldEqual, a.Markdown)
					}
					So(len(model.Breadcrumb), ShouldEqual, len(breadcrumbs))
					for i, b := range breadcrumbs {
						found := model.Breadcrumb[i]
						if i == 0 && b.Description.Title == "Home" {
							So(found.Title, ShouldEqual, "Hafan")
						} else {
							So(found.Title, ShouldEqual, b.Description.Title)
						}
						So(found.URI, ShouldEqual, b.URI)
					}
				})

				Convey("CreateSixteensBulletinModel maps correctly", func() {
					model := CreateSixteensBulletinModel(basePage, bulletin, breadcrumbs, "cy")

					So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
					So(model.FeatureFlags.SixteensVersion, ShouldEqual, "67f6982")
					So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
					So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
					So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
					So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
					So(model.Summary, ShouldEqual, bulletin.Description.Summary)
					So(model.Type, ShouldEqual, bulletin.Type)
					So(model.URI, ShouldEqual, bulletin.URI)
					So(model.ParentPath, ShouldEqual, "the/bulletin/uri/path")
					So(model.CorrectedPath, ShouldBeEmpty)
					So(model.Edition, ShouldEqual, bulletin.Description.Edition)
					So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
					So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
					So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
					So(model.LatestRelease, ShouldBeTrue)
					So(model.LatestReleaseUri, ShouldEqual, bulletin.LatestReleaseURI)
					So(model.DatasetId, ShouldEqual, bulletin.Description.DatasetID)
					So(len(model.Sections), ShouldEqual, len(bulletin.Sections))
					assertSections(model.Sections, bulletin.Sections)
					assertSections(model.Accordion, bulletin.Accordion)
					assertLinks(model.RelatedBulletins, bulletin.RelatedBulletins)
					assertLinks(model.RelatedData, bulletin.RelatedData)
					assertLinks(model.Links, bulletin.Links)
					assertFigures(model.Charts, bulletin.Charts)
					assertFigures(model.Tables, bulletin.Tables)
					assertFigures(model.Images, bulletin.Images)
					assertFigures(model.Equations, bulletin.Equations)
					So(len(model.Versions), ShouldEqual, len(bulletin.Versions))
					// Versions should be sorted by date
					sort.Slice(bulletin.Versions, func(i, j int) bool { return bulletin.Versions[i].ReleaseDate > bulletin.Versions[j].ReleaseDate })
					for i, v := range bulletin.Versions {
						f := model.Versions[i]
						So(f.Date, ShouldEqual, v.ReleaseDate)
						So(f.Markdown, ShouldEqual, v.Notice)
						So(f.URI, ShouldEqual, v.URI)
					}
					So(len(model.Alerts), ShouldEqual, len(bulletin.Alerts))
					// Alerts should be sorted by date
					sort.Slice(bulletin.Alerts, func(i, j int) bool { return bulletin.Alerts[i].Date > bulletin.Alerts[j].Date })
					for i, a := range bulletin.Alerts {
						f := model.Alerts[i]
						So(f.Date, ShouldEqual, a.Date)
						So(f.Markdown, ShouldEqual, a.Markdown)
					}
					So(len(model.Breadcrumb), ShouldEqual, len(breadcrumbs))
					for i, b := range breadcrumbs {
						found := model.Breadcrumb[i]
						So(found.Title, ShouldEqual, b.Description.Title)
						So(found.URI, ShouldEqual, b.URI)
					}
				})
			})
			Convey("When the bulletin URI is a previous version", func() {
				bulletin.URI = "the/bulletin/uri/path/previous/version"
				Convey("CreateBulletinModel maps correctly", func() {
					requestProtocol := "https"
					model := CreateBulletinModel(basePage, bulletin, breadcrumbs, "cy", requestProtocol, serviceMessage, bannerData)

					So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
					So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
					So(model.BetaBannerEnabled, ShouldBeTrue)
					So(model.ServiceMessage, ShouldEqual, serviceMessage)
					So(model.EmergencyBanner.Title, ShouldEqual, emergencyBannerTitle)
					So(model.EmergencyBanner.Type, ShouldEqual, emergencyBannerType)
					So(model.EmergencyBanner.Description, ShouldEqual, emergencyBannerDescription)
					So(model.EmergencyBanner.URI, ShouldEqual, emergencyBannerUri)
					So(model.EmergencyBanner.LinkText, ShouldEqual, emergencyBannerLinkText)
					So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
					So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
					So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
					So(model.Type, ShouldEqual, bulletin.Type)
					So(model.URI, ShouldEqual, bulletin.URI)
					So(model.ParentPath, ShouldEqual, "the/bulletin/uri/path/previous")
					So(model.CorrectedPath, ShouldEqual, "the/bulletin/uri/path")
					So(model.Edition, ShouldEqual, bulletin.Description.Edition)
					So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
					So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
					So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
					So(model.LatestRelease, ShouldBeTrue)
					So(model.LatestReleaseUri, ShouldEqual, bulletin.LatestReleaseURI)
					So(model.DatasetId, ShouldEqual, bulletin.Description.DatasetID)
					So(model.Census2021, ShouldEqual, false)
					So(model.AboutTheData, ShouldEqual, false)
					So(model.PreGTMJavaScript, ShouldNotBeEmpty)
					So(len(model.PreGTMJavaScript), ShouldEqual, 1)
					expectedPreGTMJs := createPreGTMJs(bulletin.Description.Title, bulletin.Description.ReleaseDate, bulletin.URI, "")
					So(model.PreGTMJavaScript[0], ShouldResemble, expectedPreGTMJs)
					So(len(model.Sections), ShouldEqual, len(bulletin.Sections))
					assertSections(model.Sections, bulletin.Sections)
					assertSections(model.Accordion, bulletin.Accordion)
					assertContentsView(model.ContentsView, bulletin.Sections, bulletin.Accordion, model.AboutTheData)
					assertLinks(model.RelatedBulletins, bulletin.RelatedBulletins)
					assertLinks(model.RelatedData, bulletin.RelatedData)
					assertLinks(model.Links, bulletin.Links)
					assertFigures(model.Charts, bulletin.Charts)
					assertFigures(model.Tables, bulletin.Tables)
					assertFigures(model.Images, bulletin.Images)
					assertFigures(model.Equations, bulletin.Equations)
					assertShareLinks(model.ShareLinks, bulletin.URI, requestProtocol)
					So(len(model.Versions), ShouldEqual, len(bulletin.Versions))
					// Versions should be sorted by date
					sort.Slice(bulletin.Versions, func(i, j int) bool { return bulletin.Versions[i].ReleaseDate > bulletin.Versions[j].ReleaseDate })
					for i, v := range bulletin.Versions {
						f := model.Versions[i]
						So(f.Date, ShouldEqual, v.ReleaseDate)
						So(f.Markdown, ShouldEqual, v.Notice)
						So(f.URI, ShouldEqual, v.URI)
					}
					So(len(model.Alerts), ShouldEqual, len(bulletin.Alerts))
					// Alerts should be sorted by date
					sort.Slice(bulletin.Alerts, func(i, j int) bool { return bulletin.Alerts[i].Date > bulletin.Alerts[j].Date })
					for i, a := range bulletin.Alerts {
						f := model.Alerts[i]
						So(f.Date, ShouldEqual, a.Date)
						So(f.Markdown, ShouldEqual, a.Markdown)
					}
					So(len(model.Breadcrumb), ShouldEqual, len(breadcrumbs))
					for i, b := range breadcrumbs {
						found := model.Breadcrumb[i]
						if i == 0 && b.Description.Title == "Home" {
							So(found.Title, ShouldEqual, "Hafan")
						} else {
							So(found.Title, ShouldEqual, b.Description.Title)
						}
						So(found.URI, ShouldEqual, b.URI)
					}
				})

				Convey("CreateSixteensBulletinModel maps correctly", func() {
					model := CreateSixteensBulletinModel(basePage, bulletin, breadcrumbs, "cy")

					So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
					So(model.FeatureFlags.SixteensVersion, ShouldEqual, "67f6982")
					So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
					So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
					So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
					So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
					So(model.Type, ShouldEqual, bulletin.Type)
					So(model.URI, ShouldEqual, bulletin.URI)
					So(model.ParentPath, ShouldEqual, "the/bulletin/uri/path/previous")
					So(model.CorrectedPath, ShouldEqual, "the/bulletin/uri/path")
					So(model.Edition, ShouldEqual, bulletin.Description.Edition)
					So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
					So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
					So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
					So(model.LatestRelease, ShouldBeTrue)
					So(model.LatestReleaseUri, ShouldEqual, bulletin.LatestReleaseURI)
					So(model.DatasetId, ShouldEqual, bulletin.Description.DatasetID)
					So(len(model.Sections), ShouldEqual, len(bulletin.Sections))
					assertSections(model.Sections, bulletin.Sections)
					assertSections(model.Accordion, bulletin.Accordion)
					assertLinks(model.RelatedBulletins, bulletin.RelatedBulletins)
					assertLinks(model.RelatedData, bulletin.RelatedData)
					assertLinks(model.Links, bulletin.Links)
					assertFigures(model.Charts, bulletin.Charts)
					assertFigures(model.Tables, bulletin.Tables)
					assertFigures(model.Images, bulletin.Images)
					assertFigures(model.Equations, bulletin.Equations)
					So(len(model.Versions), ShouldEqual, len(bulletin.Versions))
					// Versions should be sorted by date
					sort.Slice(bulletin.Versions, func(i, j int) bool { return bulletin.Versions[i].ReleaseDate > bulletin.Versions[j].ReleaseDate })
					for i, v := range bulletin.Versions {
						f := model.Versions[i]
						So(f.Date, ShouldEqual, v.ReleaseDate)
						So(f.Markdown, ShouldEqual, v.Notice)
						So(f.URI, ShouldEqual, v.URI)
					}
					So(len(model.Alerts), ShouldEqual, len(bulletin.Alerts))
					// Alerts should be sorted by date
					sort.Slice(bulletin.Alerts, func(i, j int) bool { return bulletin.Alerts[i].Date > bulletin.Alerts[j].Date })
					for i, a := range bulletin.Alerts {
						f := model.Alerts[i]
						So(f.Date, ShouldEqual, a.Date)
						So(f.Markdown, ShouldEqual, a.Markdown)
					}
					So(len(model.Breadcrumb), ShouldEqual, len(breadcrumbs))
					for i, b := range breadcrumbs {
						found := model.Breadcrumb[i]
						So(found.Title, ShouldEqual, b.Description.Title)
						So(found.URI, ShouldEqual, b.URI)
					}
				})
			})
		})
	})
}

func assertShareLinks(shareLinks ShareLinks, uri, requestProtocol string) {
	So(shareLinks, ShouldContainKey, model.SocialEmail.String())
	So(shareLinks, ShouldContainKey, model.SocialLinkedin.String())
	So(shareLinks, ShouldContainKey, model.SocialTwitter.String())

	emailUrl, err := url.Parse(shareLinks[model.SocialEmail.String()].Url)
	So(err, ShouldBeNil)
	emailParams := emailUrl.Query()
	So(emailParams.Get("body"), ShouldContainSubstring, requestProtocol)
	So(emailParams.Get("body"), ShouldContainSubstring, uri)

	linkedinUrl, err := url.Parse(shareLinks[model.SocialLinkedin.String()].Url)
	So(err, ShouldBeNil)
	linkedinParams := linkedinUrl.Query()
	So(linkedinParams.Get("url"), ShouldContainSubstring, requestProtocol)
	So(linkedinParams.Get("url"), ShouldContainSubstring, uri)

	twitterUrl, err := url.Parse(shareLinks[model.SocialTwitter.String()].Url)
	So(err, ShouldBeNil)
	twitterParams := twitterUrl.Query()
	So(twitterParams.Get("url"), ShouldContainSubstring, requestProtocol)
	So(twitterParams.Get("url"), ShouldContainSubstring, uri)
}

func assertContentsView(found []ViewSection, expectedSections, expectedAccordion []zebedee.Section, aboutTheData bool) {
	totalSections := len(expectedSections)
	totalAccordions := len(expectedAccordion)
	expectedSectionCount := totalSections + totalAccordions
	if aboutTheData {
		expectedSectionCount++
	}
	So(len(found), ShouldEqual, expectedSectionCount)
	for i := range expectedSections {
		So(found[i].Type, ShouldEqual, "section")
	}
	for i := range expectedAccordion {
		So(found[totalAccordions+i].Type, ShouldEqual, "accordion")
	}
}

func assertSections(found []Section, expected []zebedee.Section) {
	So(len(found), ShouldEqual, len(expected))
	for i, s := range expected {
		So(found[i].Title, ShouldEqual, s.Title)
		So(found[i].Markdown, ShouldEqual, s.Markdown)
	}
}

func assertLinks(found []Link, expected []zebedee.Link) {
	So(len(found), ShouldEqual, len(expected))
	for i, s := range expected {
		So(found[i].Title, ShouldEqual, s.Title)
		So(found[i].URI, ShouldEqual, s.URI)
	}
}

func assertFigures(found []Figure, expected []zebedee.Figure) {
	So(len(found), ShouldEqual, len(expected))
	for i, s := range expected {
		So(found[i].Title, ShouldEqual, s.Title)
		So(found[i].URI, ShouldEqual, s.URI)
		So(found[i].Filename, ShouldEqual, s.Filename)
		So(found[i].Version, ShouldEqual, s.Version)
	}
}

func createPreGTMJs(title, releaseDate, url, tag string) template.JS {
	return template.JS("dataLayer.push({" +
		"\"analyticsOptOut\": getUsageCookieValue()," +
		"\"gtm.whitelist\": [\"google\",\"hjtc\",\"lcl\"]," +
		"\"gtm.blacklist\": [\"customScripts\",\"sp\",\"adm\",\"awct\",\"k\",\"d\",\"j\"]," +
		"\"contentTitle\": \"" + title + "\"," +
		"\"release-date-status\": \"" + releaseDate + "\"," +
		"\"url\": \"" + url + "\"," +
		"\"tag\": \"" + tag + "\"" +
		"});")
}
