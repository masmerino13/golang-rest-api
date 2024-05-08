package controllers

import (
	"fmt"
	"net/http"

	"lens.com/m/v2/models"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	user, err := u.UserService.Create(models.NewUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})

	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("error creating user: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "user created: %+v", user)
}
