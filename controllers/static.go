package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
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

func FAQ(tpl Template) http.HandlerFunc {
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
