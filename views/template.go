package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"lens.com/m/v2/helpers"
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

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return `<input type="hidden">`
		},
		"isLogin": func() bool {
			return false
		},
	})

	tpl, err := tpl.ParseFS(fs, patterns...)

	if err != nil {
		return Template{}, fmt.Errorf("error parsing %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := t.htmlTpl.Clone()

	if err != nil {
		log.Printf("Error cloning template %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"isLogin": func() bool {
				cookie, err := r.Cookie(helpers.CookieAuthToken)

				if err != nil {
					log.Printf("Error reading cookie %v", err)
					return false
				}

				err = cookie.Valid()

				return err == nil
			},
		},
	)
	w.Header().Set("Content-Type", "text/html")

	// NOTE: will cause performance issues if there are many HTML pages with a lot of HTML content
	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)

	if err != nil {
		log.Printf("Error executing templae %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

func Msg(m string) (string, string) {
	return m, m
}
