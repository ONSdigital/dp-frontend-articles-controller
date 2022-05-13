package mapper

import (
	"sort"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/articles"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	coreModel "github.com/ONSdigital/dp-renderer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {

	Convey("Given a bulletin, basePage and breadcrumbs", t, func() {
		basePage := coreModel.NewPage("path/to/assets", "site-domain")

		bulletin := articles.Bulletin{
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
		Convey("When the bulletin URI is not a previous version", func() {
			bulletin.URI = "the/bulletin/uri/path/version"
			Convey("CreateSixteensBulletinModel maps correctly", func() {
				model := CreateSixteensBulletinModel(basePage, bulletin, breadcrumbs)

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
			Convey("CreateSixteensBulletinModel maps correctly", func() {
				model := CreateSixteensBulletinModel(basePage, bulletin, breadcrumbs)

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

	Convey("Given a bulletin, basePage and breadcrumbs", t, func() {
		basePage := coreModel.NewPage("path/to/assets", "site-domain")

		bulletin := articles.Bulletin{
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
		Convey("When the bulletin URI is not a previous version", func() {
			bulletin.URI = "the/bulletin/uri/path/version"
			Convey("CreateBulletinModel maps correctly", func() {
				model := CreateBulletinModel(basePage, bulletin, breadcrumbs)

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
				model := CreateBulletinModel(basePage, bulletin, breadcrumbs)

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
