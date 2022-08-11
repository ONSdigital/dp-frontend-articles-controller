/*
{
  "downloads": [
    {
      "title": "ANWAR",
      "file": "deccanwar2november2015_tcm77-430106.xls"
    },
    {
      "title": "Tables",
      "file": "decctablesnovember2015_tcm77-430108.pdf"
    }
  ],
  "markdown": [
    "Export and import estimates of crude oil and other fuels including adjustments. "
  ],
  "links": [],
  "type": "static_adhoc",
  "uri": "/economy/environmentalaccounts/adhocs/005204fuelandenergydataprovidedonamonthlybasistodecc",
  "description": {
    "title": "Fuel and energy data provided on a monthly basis to DECC",
    "keywords": [],
    "metaDescription": "Export and import estimates of crude oil and other fuels including adjustments. ",
    "releaseDate": "2016-01-15T08:03:43.411Z",
    "unit": "",
    "preUnit": "",
    "source": "",
    "reference": "005204"
  }
}
*/

package model

import (
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type DownloadJSON struct {
	Title string `json:"title"`
	File  string `json:"file"`
}

type LinkJSON struct {
}

type DescriptionJSON struct {
	Title           string   `json:"title"`
	Keywords        []string `json:"keywords"`
	MetaDescription string   `json:"metaDescription"`
	ReleaseDate     string   `json:"releaseDate"`
	Unit            string   `json:"unit"`
	PreUnit         string   `json:"preUnit"`
	Source          string   `json:"source"`
	Reference       string   `json:"reference"`
}

type AdHocJSON struct {
	Downloads   []DownloadJSON  `json:"downloads"`
	Markdown    []string        `json:"markdown"`
	Links       []LinkJSON      `json:"links"`
	Type        string          `json:"type"`
	Uri         string          `json:"uri"`
	Description DescriptionJSON `json:"description"`
}

type Download struct {
	Title string `json:"title"`
	File  string `json:"file"`
}

type AdHoc struct {
	coreModel.Page
	BodyMarkdown []string   `json:"body_markdown"`
	Reference    string     `json:"reference"`
	Downloads    []Download `json:"downloads"`
}
