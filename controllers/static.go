package controllers

import (
	"html/template"
	"net/http"

	"lens.com/m/v2/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

type Question struct {
	Question string
	Answer   template.HTML
}

type FAQData struct {
	Questions []Question
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []Question{
		{
			Question: "You hungry?",
			Answer:   "Yes, I'm",
		},
		{
			Question: "You hungry?",
			Answer:   "Yes, I'm",
		},
		{
			Question: "You hungry?",
			Answer:   "Yes, I'm",
		},
		{
			Question: "You happy?",
			Answer:   "no, I'm not",
		},
	}

	data := FAQData{
		Questions: questions,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, data)
	}
}
