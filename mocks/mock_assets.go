package mocks

import "strings"

var cyLocale = []string{
	"[BreadcrumbHome]",
	"one=\"Hafan\"",
}

var enLocale = []string{
	"[BreadcrumbHome]",
	"one=\"Home\"",
}

func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte(strings.Join(cyLocale, "\n")), nil
	}
	return []byte(strings.Join(enLocale, "\n")), nil
}
