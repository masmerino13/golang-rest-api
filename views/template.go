package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}

	return t
}

func ParseFS(fs fs.FS, patterns string) (Template, error) {
	tpl, err := template.ParseFS(fs, patterns)

	if err != nil {
		return Template{}, fmt.Errorf("error parsing %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, nil
}

func Parse(filePath string) (Template, error) {
	tpl, err := template.ParseFiles(filePath)

	if err != nil {
		return Template{}, fmt.Errorf("error parsing %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	err := t.htmlTpl.Execute(w, nil)

	if err != nil {
		log.Printf("Error executing templae %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func Msg(m string) (string, string) {
	return m, m
}
