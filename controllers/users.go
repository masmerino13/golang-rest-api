package controllers

import (
	"fmt"
	"net/http"

	"lens.com/m/v2/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	u.Templates.SignIn.Execute(w, r, nil)
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

	session, err := u.SessionService.Create(user.ID)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "signin", http.StatusFound)
		return
	}

	cookie := http.Cookie{
		Name:     "authToken",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/users/me", http.StatusFound)

	fmt.Fprintf(w, "user created: %+v", user)
}

func (u Users) Auth(w http.ResponseWriter, r *http.Request) {
	user, err := u.UserService.Authenticate(r.FormValue("email"), r.FormValue("password"))

	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("error authenticating user: %s", err), http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)

	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("error authenticating user: %s", err), http.StatusUnauthorized)
		return
	}

	cookie := http.Cookie{
		Name:     "authToken",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "user authenticated: %+v", user)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	authToken, err := r.Cookie("authUser")

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	user, err := u.SessionService.User(authToken.Value)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "current user: %s", user.Email)
}
