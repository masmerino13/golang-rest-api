package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"lens.com/m/v2/views"
)

func execTemplate(w http.ResponseWriter, filePath string) {
	w.Header().Set("Content-Type", "text/html")

	tpl, err := views.Parse(filePath)

	if err != nil {
		log.Printf("Error parsin templae %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	fmt.Println(views.Msg("el mensaje"))

	tpl.Execute(w, nil)
}

func homeFunc(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")

	execTemplate(w, tplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")

	execTemplate(w, tplPath)
}

type Router struct{}

func (Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch path {
	case "/":
		homeFunc(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func main() {
	r := chi.NewRouter()

	port := "3002"

	fmt.Println("Starting server at ", port)

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", homeFunc)
	r.Get("/contact", contactHandler)
	r.Route("/products", func(r chi.Router) {
		r.Get("/{productID}", func(w http.ResponseWriter, r *http.Request) {
			productID := chi.URLParam(r, "productID")

			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, "<h1>Product here</h1>")
			fmt.Fprint(w, productID)
		})
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found blah", http.StatusNotFound)
	})

	http.ListenAndServe(":"+port, r)
}
