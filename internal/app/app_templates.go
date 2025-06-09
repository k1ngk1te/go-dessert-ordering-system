package app

import "html/template"

func NewApplicationTemplates() (*template.Template, error) {
	var templates *template.Template // Initiate Template

	templates, err := template.ParseFiles("./templates/index.html", "./templates/login.html", "./templates/register.html")
	if err != nil {
		return nil, err
	}

	return templates, nil
}