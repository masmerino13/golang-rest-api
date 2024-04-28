package main

import (
	"fmt"
	"net/http"
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
	port := "3002"

	var router Router

	fmt.Println("Starting server at ", port)

	http.ListenAndServe(":"+port, router)
}
