package main

import (
	"embed"
	"html/template"
)

var templates *template.Template

//go:embed templates/*.htm
var templatesFS embed.FS

func initTepmplates() error {
	var err error

	templates, err = template.ParseFS(templatesFS, "templates/*.htm")
	if err != nil {
		return err
	}

	templatesFS = embed.FS{}

	return nil
}
