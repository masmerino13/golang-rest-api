package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Hi Ricardo it works!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>contact page</h1>")
}

type Router struct{}

func (Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch path {
	case "/":
		handleFunc(w, r)
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

	r.Get("/", handleFunc)
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
