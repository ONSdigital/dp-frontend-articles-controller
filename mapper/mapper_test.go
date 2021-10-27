package mapper

import (
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	coreModel "github.com/ONSdigital/dp-renderer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	basePage := coreModel.NewPage("path/to/assets", "site-domain")
	bulletin := zebedee.Bulletin{
		Type: "bulletin",
		URI:  "uri",
		Description: zebedee.Description{
			Title:             "Title",
			Edition:           "2021",
			Keywords:          []string{"energy", "waste"},
			MetaDescription:   "description",
			NationalStatistic: true,
			Contact: zebedee.Contact{
				Name:      "contact name",
				Email:     "contact@ons.gov.uk",
				Telephone: "+44",
			},
			ReleaseDate: "2015-07-08T23:00:00.000Z",
			NextRelease: "",
			DatasetID:   "22",
		},
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
	}
	Convey("CreateBulletinModel maps correctly", t, func() {
		model := CreateBulletinModel(basePage, bulletin)

		So(model.Page.PatternLibraryAssetsPath, ShouldEqual, basePage.PatternLibraryAssetsPath)
		So(model.Page.SiteDomain, ShouldEqual, basePage.SiteDomain)
		So(model.Metadata.Title, ShouldEqual, bulletin.Description.Title)
		So(model.Metadata.Description, ShouldEqual, bulletin.Description.MetaDescription)
		So(model.Metadata.Keywords, ShouldResemble, bulletin.Description.Keywords)
		So(model.Type, ShouldEqual, bulletin.Type)
		So(model.URI, ShouldEqual, bulletin.URI)
		So(model.Edition, ShouldEqual, bulletin.Description.Edition)
		So(model.NationalStatistic, ShouldEqual, bulletin.Description.NationalStatistic)
		So(model.ReleaseDate, ShouldEqual, bulletin.Description.ReleaseDate)
		So(model.NextRelease, ShouldEqual, bulletin.Description.NextRelease)
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
