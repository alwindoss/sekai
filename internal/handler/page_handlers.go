package handler

import (
	"net/http"
	"sekai/data"
	"sekai/templates/layouts"
	"sekai/templates/pages"
	"time"

	"github.com/a-h/templ"
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

func (ph *PageHandler) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	// var creds Credentials

	userName := r.FormValue("username")
	password := r.FormValue("password")

	// Get the JSON body and decode into credentials
	// err := json.NewDecoder(r.Body).Decode(&creds)
	// if err != nil {
	// 	// If the structure of the body is wrong, return an HTTP error
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// Get the expected password from our in memory map
	expectedPassword, ok := users[userName]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(15 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(ph.JWTSecret)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure:   true,
	})

	// http.RedirectHandler("/pages/welcome", http.StatusPermanentRedirect)
	http.Redirect(w, r, "/pages/welcome", http.StatusSeeOther)
}

func (ph *PageHandler) WelcomePageHandler(w http.ResponseWriter, r *http.Request) {
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
	// w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	services := []pages.AppService{
		{
			Name:        "Car Management",
			Description: "Maintenance logs, fuel tracking, and insurance.",
			Icon:        "🚗",
			Link:        "/cars",
			Status:      "Up to date",
		},
		{
			Name:        "Finance",
			Description: "Monthly budget and investment tracking.",
			Icon:        "💰",
			Link:        "/finance",
			Status:      "Review needed",
		},
		{
			Name:        "Reminders",
			Description: "Filter changes, trash day, and meds.",
			Icon:        "🔔",
			Link:        "/reminders",
			Status:      "3 active",
		},
		{
			Name:        "Smart Home",
			Description: "Control lights, locks, and thermostats.",
			Icon:        "🏠",
			Link:        "/home-iot",
			Status:      "Online",
		},
	}
	templ.Handler(layouts.Default(pages.HomePage(services), &data.PageData{
		Title: "About",
	})).ServeHTTP(w, r)

}
