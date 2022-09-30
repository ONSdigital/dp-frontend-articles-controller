package mapper

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/articles"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-renderer/helper"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type ShareLinks map[string]coreModel.ShareLink

type BulletinModel struct {
	coreModel.Page
	Summary           string        `json:"summary"`
	Sections          []Section     `json:"sections"`
	Accordion         []Section     `json:"accordion"`
	Charts            []Figure      `json:"charts"`
	Tables            []Figure      `json:"tables"`
	Images            []Figure      `json:"images"`
	Equations         []Figure      `json:"equations"`
	RelatedBulletins  []Link        `json:"relatedBulletins"`
	RelatedData       []Link        `json:"relatedData"`
	Links             []Link        `json:"links"`
	URI               string        `json:"uri"`
	NationalStatistic bool          `json:"nationalStatistic"`
	LatestRelease     bool          `json:"latestRelease"`
	Edition           string        `json:"edition"`
	ReleaseDate       string        `json:"releaseDate"`
	NextRelease       string        `json:"nextRelease"`
	Contact           Contact       `json:"contact"`
	Versions          []Message     `json:"versions"`
	Alerts            []Message     `json:"alerts"`
	ParentPath        string        `json:"parentPath"`
	CorrectedPath     string        `json:"correctedPath"`
	LatestReleaseUri  string        `json:"latestReleaseUri"`
	ContentsView      []ViewSection `json:"contentsView"`
	ShareLinks        ShareLinks    `json:"shareLinks"`
	Census2021        bool          `json:"census_2021"`
	AboutTheData      bool          `json:"about_the_data"`
	Auxiliary         []Section     `json:"auxiliary"`
}

// Intermediate view to aid template rendering of Sections and Accordion
type ViewSection struct {
	Id       string
	Type     string
	Source   *[]Section
	Index    int
	BackTo   coreModel.BackTo
	Language string
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

func mapEmergencyBanner(bannerData zebedee.EmergencyBanner) coreModel.EmergencyBanner {
	var mappedEmergencyBanner coreModel.EmergencyBanner
	emptyBannerObj := zebedee.EmergencyBanner{}
	if bannerData != emptyBannerObj {
		mappedEmergencyBanner.Title = bannerData.Title
		mappedEmergencyBanner.Type = strings.Replace(bannerData.Type, "_", "-", -1)
		mappedEmergencyBanner.Description = bannerData.Description
		mappedEmergencyBanner.URI = bannerData.URI
		mappedEmergencyBanner.LinkText = bannerData.LinkText
	}
	return mappedEmergencyBanner
}

func CreateBulletinModel(basePage coreModel.Page, bulletin articles.Bulletin, bcs []zebedee.Breadcrumb, lang, requestProtocol, serviceMessage string, emergencyBannerContent zebedee.EmergencyBanner) BulletinModel {
	model := BulletinModel{
		Page: basePage,
	}
	model.Language = lang
	model.ServiceMessage = serviceMessage
	model.EmergencyBanner = mapEmergencyBanner(emergencyBannerContent)
	model.BetaBannerEnabled = true
	model.Type = bulletin.Type
	model.Census2021 = bulletin.Description.Survey == "census"
	model.AboutTheData = model.Census2021
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

	model.Page.Breadcrumb = mapBreadcrumbTrail(bcs, model.Language)
	populateContents(&model)

	currentUrl := getCurrentUrl(requestProtocol, model.SiteDomain, model.URI, lang)
	model.ShareLinks = createShareLinks(model.Metadata.Title, currentUrl)

	return model
}

// Sections are always followed by Accordions
func populateContents(model *BulletinModel) {
	appendSections := func(list *[]Section, listType string, views *[]ViewSection) {
		for index := range *list {
			*views = append(*views, ViewSection{
				Id:     fmt.Sprintf("%s-%d", listType, index),
				Type:   listType,
				Source: list,
				Index:  index,
			})
		}
	}

	views := make([]ViewSection, 0, len(model.Sections)+len(model.Accordion))
	appendSections(&model.Sections, "section", &views)
	appendSections(&model.Accordion, "accordion", &views)

	for index := range views {
		views[index].BackTo = coreModel.BackTo{
			Text: coreModel.Localisation{
				LocaleKey: "BackToContents",
				Plural:    4,
			},
			AnchorFragment: "toc",
		}
		views[index].Language = model.Language
	}

	if model.AboutTheData {
		model.Auxiliary = append(model.Auxiliary, Section{
			Title:    helper.Localise("PageSectionAboutTheData", model.Language, 1),
			Markdown: helper.Localise("AboutTheDataMarkdown", model.Language, 1),
		})
		views = append(views, ViewSection{
			Id:     "aboutthedata",
			Type:   "auxiliary",
			Source: &model.Auxiliary,
			Index:  0,
		})
		views[len(views)-1].BackTo = coreModel.BackTo{
			Text: coreModel.Localisation{
				LocaleKey: "BackToContents",
				Plural:    4,
			},
			AnchorFragment: "toc",
		}
		views[len(views)-1].Language = model.Language
	}

	model.TableOfContents = createTableOfContents(views)
	model.ContentsView = views
}

func createShareLinks(title, url string) ShareLinks {
	return ShareLinks{
		coreModel.SocialEmail.String():    coreModel.SocialEmail.CreateLink(title, url),
		coreModel.SocialLinkedin.String(): coreModel.SocialLinkedin.CreateLink(title, url),
		coreModel.SocialTwitter.String():  coreModel.SocialTwitter.CreateLink(title, url),
	}
}

func getCurrentUrl(requestProtocol, siteDomain, path, lang string) string {
	var subDomain string
	if lang == "cy" {
		subDomain = "cy."
	}

	if siteDomain == "localhost" || siteDomain == "" {
		siteDomain = "ons.gov.uk"
	}

	currentUrl := url.URL{
		Scheme: requestProtocol,
		Host:   subDomain + siteDomain,
		Path:   path,
	}

	return currentUrl.String()
}

func createTableOfContents(views []ViewSection) coreModel.TableOfContents {
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
	displayOrder := make([]string, 0, len(views))

	for _, view := range views {
		section := (*view.Source)[view.Index]
		sections[view.Id] = coreModel.ContentSection{
			Current: false,
			Title: coreModel.Localisation{
				Text: section.Title,
			},
		}
		displayOrder = append(displayOrder, view.Id)
	}

	toc.Sections = sections
	toc.DisplayOrder = displayOrder

	return toc
}

func parentPath(p string) string {
	return p[:strings.LastIndex(p, "/")]
}

type SectionReference struct {
	Type  string
	Index int
}

func mapBreadcrumbTrail(crumbs []zebedee.Breadcrumb, language string) []coreModel.TaxonomyNode {
	trail := []coreModel.TaxonomyNode{}

	for _, crumb := range crumbs {
		trail = append(trail, coreModel.TaxonomyNode{
			Title: crumb.Description.Title,
			URI:   crumb.URI,
		})
	}

	if len(trail) > 0 && trail[0].Title == "Home" {
		trail[0].Title = helper.Localise("BreadcrumbHome", language, 1)
	}

	return trail
}
