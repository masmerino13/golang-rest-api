package helpers

import "net/http"

const (
	CookieAuthToken = "authToken"
)

func NewCookie(name, value string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}

	return &cookie
}

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := NewCookie(name, value)

	http.SetCookie(w, cookie)
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := NewCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
