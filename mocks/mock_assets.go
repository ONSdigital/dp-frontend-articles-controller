package mocks

import "strings"

var cyLocale = []string{
	"[BreadcrumbHome]",
	"one=\"Hafan\"",
	"[PageSectionAboutTheData]",
	"one=\"Ynglŷn â'r data\"",
	"[AboutTheDataMarkdown]",
	"one=\"This release includes data from Census 2021 (cy)\"",
}

var enLocale = []string{
	"[BreadcrumbHome]",
	"one=\"Home\"",
	"[PageSectionAboutTheData]",
	"one=\"About the data\"",
	"[AboutTheDataMarkdown]",
	"one=\"This release includes data from Census 2021\"",
}

func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte(strings.Join(cyLocale, "\n")), nil
	}
	return []byte(strings.Join(enLocale, "\n")), nil
}
