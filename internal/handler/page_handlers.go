package handler

import (
	"fmt"
	"net/http"
	"sekai/data"
	"sekai/templates/layouts"
	"sekai/templates/pages"

	"github.com/golang-jwt/jwt/v5"
)

type PageHandler struct {
	JWTSecret []byte
}

func (ph *PageHandler) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	layouts.Default(pages.Index(), &data.PageData{
		Title: "Sekai | Home",
	}).Render(r.Context(), w)
}

func (ph *PageHandler) AboutPageHandler(w http.ResponseWriter, r *http.Request) {
	layouts.Default(pages.About(), &data.PageData{
		Title: "Sekai | About",
	}).Render(r.Context(), w)
}

func (ph *PageHandler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	layouts.Default(pages.LoginPage(), &data.PageData{
		Title: "Sekai | Login",
	}).Render(r.Context(), w)
	// templ.Handler().ServeHTTP
}

func (ph *PageHandler) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return ph.JWTSecret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}
