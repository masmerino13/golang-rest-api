package controllers

import (
	"fmt"
	"net/http"

	"lens.com/m/v2/context"
	"lens.com/m/v2/helpers"
	"lens.com/m/v2/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	// NOTE: Here we define the model, BUT it's required to be set in main.go in order to work
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

	helpers.SetCookie(w, helpers.CookieAuthToken, session.Token)

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

	helpers.SetCookie(w, helpers.CookieAuthToken, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	fmt.Fprintf(w, "current user: %s", user.Email)
}

func (u Users) SingOut(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("token")
	authToken, err := r.Cookie(helpers.CookieAuthToken)

	fmt.Printf("token: %v", authToken)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(authToken.Value)

	fmt.Printf("token 2: %v", err)

	if err != nil {
		fmt.Fprintf(w, "Error singin out")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Printf("pasa")

	helpers.DeleteCookie(w, helpers.CookieAuthToken)

	http.Redirect(w, r, "/signin", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := r.Cookie(helpers.CookieAuthToken)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(authToken.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithUser(r.Context(), user)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())

		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
