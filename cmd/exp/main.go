package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")

	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Ric Cena",
	}

	err = t.Execute(os.Stdout, user)

	if err != nil {
		panic(err)
	}
}
