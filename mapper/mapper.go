package mapper

import (
	"sort"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/articles"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type BulletinModel struct {
	coreModel.Page
	Summary           string    `json:"summary"`
	Sections          []Section `json:"sections"`
	Accordion         []Section `json:"accordion"`
	Charts            []Figure  `json:"charts"`
	Tables            []Figure  `json:"tables"`
	Images            []Figure  `json:"images"`
	Equations         []Figure  `json:"equations"`
	RelatedBulletins  []Link    `json:"relatedBulletins"`
	RelatedData       []Link    `json:"relatedData"`
	Links             []Link    `json:"links"`
	URI               string    `json:"uri"`
	NationalStatistic bool      `json:"nationalStatistic"`
	LatestRelease     bool      `json:"latestRelease"`
	Edition           string    `json:"edition"`
	ReleaseDate       string    `json:"releaseDate"`
	NextRelease       string    `json:"nextRelease"`
	Contact           Contact   `json:"contact"`
	Versions          []Message `json:"versions"`
	Alerts            []Message `json:"alerts"`
	ParentPath        string    `json:"parentPath"`
	CorrectedPath     string    `json:"correctedPath"`
	LatestReleaseUri  string    `json:"latestReleaseUri"`
}

type Contact struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

type Link struct {
	Title string `json:"title"`
	URI   string `json:"uri"`
}

type Figure struct {
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
	URI      string `json:"uri"`
}

type Section struct {
	Title    string `json:"title"`
	Markdown string `json:"markdown"`
}

type Message struct {
	Date     string `json:"date"`
	Markdown string `json:"markdown"`
	URI      string `json:"uri"`
}

func CreateSixteensBulletinModel(basePage coreModel.Page, bulletin articles.Bulletin, bcs []zebedee.Breadcrumb, lang string) BulletinModel {
	model := BulletinModel{
		Page: basePage,
	}
	model.Language = lang
	model.FeatureFlags.SixteensVersion = "67f6982"

	model.Metadata = coreModel.Metadata{
		Title:       bulletin.Description.Title,
		Description: bulletin.Description.MetaDescription,
		Keywords:    bulletin.Description.Keywords,
	}
	model.URI = bulletin.URI
	model.Summary = bulletin.Description.Summary
	model.DatasetId = bulletin.Description.DatasetID
	model.NationalStatistic = bulletin.Description.NationalStatistic
	model.Edition = bulletin.Description.Edition
	model.ReleaseDate = bulletin.Description.ReleaseDate
	model.NextRelease = bulletin.Description.NextRelease
	model.LatestRelease = bulletin.Description.LatestRelease
	model.LatestReleaseUri = bulletin.LatestReleaseURI
	model.Contact = Contact{
		Name:      bulletin.Description.Contact.Name,
		Email:     bulletin.Description.Contact.Email,
		Telephone: bulletin.Description.Contact.Telephone,
	}

	model.ParentPath = parentPath(bulletin.URI)
	if strings.HasSuffix(model.ParentPath, "previous") {
		model.CorrectedPath = parentPath(model.ParentPath)
	}

	model.Sections = []Section{}
	for _, s := range bulletin.Sections {
		model.Sections = append(model.Sections, Section{
			Title:    s.Title,
			Markdown: s.Markdown,
		})
	}
	model.Accordion = []Section{}
	for _, s := range bulletin.Accordion {
		model.Accordion = append(model.Accordion, Section{
			Title:    s.Title,
			Markdown: s.Markdown,
		})
	}
	model.RelatedBulletins = []Link{}
	for _, s := range bulletin.RelatedBulletins {
		model.RelatedBulletins = append(model.RelatedBulletins, Link{
			Title: s.Title,
			URI:   s.URI,
		})
	}
	model.RelatedData = []Link{}
	for _, s := range bulletin.RelatedData {
		model.RelatedData = append(model.RelatedData, Link{
			Title: s.Title,
			URI:   s.URI,
		})
	}
	model.Links = []Link{}
	for _, s := range bulletin.Links {
		model.Links = append(model.Links, Link{
			Title: s.Title,
			URI:   s.URI,
		})
	}
	model.Charts = []Figure{}
	for _, s := range bulletin.Charts {
		model.Charts = append(model.Charts, Figure{
			Title:    s.Title,
			Filename: s.Filename,
			Version:  s.Version,
			URI:      s.URI,
		})
	}
	model.Tables = []Figure{}
	for _, s := range bulletin.Tables {
		model.Tables = append(model.Tables, Figure{
			Title:    s.Title,
			Filename: s.Filename,
			Version:  s.Version,
			URI:      s.URI,
		})
	}
	model.Images = []Figure{}
	for _, s := range bulletin.Images {
		model.Images = append(model.Images, Figure{
			Title:    s.Title,
			Filename: s.Filename,
			Version:  s.Version,
			URI:      s.URI,
		})
	}
	model.Equations = []Figure{}
	for _, s := range bulletin.Equations {
		model.Equations = append(model.Equations, Figure{
			Title:    s.Title,
			Filename: s.Filename,
			Version:  s.Version,
			URI:      s.URI,
		})
	}

	model.Versions = []Message{}
	for _, v := range bulletin.Versions {
		model.Versions = append(model.Versions, Message{
			Date:     v.ReleaseDate,
			Markdown: v.Notice,
			URI:      v.URI,
		})
	}
	sort.Slice(model.Versions, func(i, j int) bool { return model.Versions[i].Date > model.Versions[j].Date })

	model.Alerts = []Message{}
	for _, a := range bulletin.Alerts {
		model.Alerts = append(model.Alerts, Message{
			Markdown: a.Markdown,
			Date:     a.Date,
		})
	}
	sort.Slice(model.Alerts, func(i, j int) bool { return model.Alerts[i].Date > model.Alerts[j].Date })

	for _, bc := range bcs {
		model.Page.Breadcrumb = append(model.Page.Breadcrumb, coreModel.TaxonomyNode{
			Title: bc.Description.Title,
			URI:   bc.URI,
		})
	}

	return model
}

func CreateBulletinModel(basePage coreModel.Page, bulletin articles.Bulletin, bcs []zebedee.Breadcrumb, lang string) BulletinModel {
	model := BulletinModel{
		Page: basePage,
	}
	model.Language = lang
	model.FeatureFlags.SixteensVersion = "67f6982"

	model.Metadata = coreModel.Metadata{
		Title:       bulletin.Description.Title,
		Description: bulletin.Description.MetaDescription,
		Keywords:    bulletin.Description.Keywords,
	}
	model.URI = bulletin.URI
	model.Summary = bulletin.Description.Summary
	model.DatasetId = bulletin.Description.DatasetID
	model.NationalStatistic = bulletin.Description.NationalStatistic
	model.Edition = bulletin.Description.Edition
	model.ReleaseDate = bulletin.Description.ReleaseDate
	model.NextRelease = bulletin.Description.NextRelease
	model.LatestRelease = bulletin.Description.LatestRelease
	model.LatestReleaseUri = bulletin.LatestReleaseURI
	model.Contact = Contact{
		Name:      bulletin.Description.Contact.Name,
		Email:     bulletin.Description.Contact.Email,
		Telephone: bulletin.Description.Contact.Telephone,
	}

	model.ParentPath = parentPath(bulletin.URI)
	if strings.HasSuffix(model.ParentPath, "previous") {
		model.CorrectedPath = parentPath(model.ParentPath)
	}

	model.Sections = []Section{}
	for _, s := range bulletin.Sections {
		model.Sections = append(model.Sections, Section{
			Title:    s.Title,
			Markdown: s.Markdown,
		})
	}
	model.Accordion = []Section{}
	for _, s := range bulletin.Accordion {
		model.Accordion = append(model.Accordion, Section{
			Title:    s.Title,
			Markdown: s.Markdown,
		})
	}
	model.RelatedBulletins = []Link{}
	for _, s := range bulletin.RelatedBulletins {
		model.RelatedBulletins = append(model.RelatedBulletins, Link{
			Title: s.Title,
			URI:   s.URI,
		})
	}
	model.RelatedData = []Link{}
	for _, s := range bulletin.RelatedData {
		model.RelatedData = append(model.RelatedData, Link{
			Title: s.Title,
			URI:   s.URI,
		})
	}
	model.Links = []Link{}
	for _, s := range bulletin.Links {
		model.Links = append(model.Links, Link{
			Title: s.Title,
			URI:   s.URI,
		})
	}
	model.Charts = []Figure{}
	for _, s := range bulletin.Charts {
		model.Charts = append(model.Charts, Figure{
			Title:    s.Title,
			Filename: s.Filename,
			Version:  s.Version,
			URI:      s.URI,
		})
	}
	model.Tables = []Figure{}
	for _, s := range bulletin.Tables {
		model.Tables = append(model.Tables, Figure{
			Title:    s.Title,
			Filename: s.Filename,
			Version:  s.Version,
			URI:      s.URI,
		})
	}
	model.Images = []Figure{}
	for _, s := range bulletin.Images {
		model.Images = append(model.Images, Figure{
			Title:    s.Title,
			Filename: s.Filename,
			Version:  s.Version,
			URI:      s.URI,
		})
	}
	model.Equations = []Figure{}
	for _, s := range bulletin.Equations {
		model.Equations = append(model.Equations, Figure{
			Title:    s.Title,
			Filename: s.Filename,
			Version:  s.Version,
			URI:      s.URI,
		})
	}

	model.Versions = []Message{}
	for _, v := range bulletin.Versions {
		model.Versions = append(model.Versions, Message{
			Date:     v.ReleaseDate,
			Markdown: v.Notice,
			URI:      v.URI,
		})
	}
	sort.Slice(model.Versions, func(i, j int) bool { return model.Versions[i].Date > model.Versions[j].Date })

	model.Alerts = []Message{}
	for _, a := range bulletin.Alerts {
		model.Alerts = append(model.Alerts, Message{
			Markdown: a.Markdown,
			Date:     a.Date,
		})
	}
	sort.Slice(model.Alerts, func(i, j int) bool { return model.Alerts[i].Date > model.Alerts[j].Date })

	for _, bc := range bcs {
		model.Page.Breadcrumb = append(model.Page.Breadcrumb, coreModel.TaxonomyNode{
			Title: bc.Description.Title,
			URI:   bc.URI,
		})
	}

	return model
}

func parentPath(p string) string {
	return p[:strings.LastIndex(p, "/")]
}
