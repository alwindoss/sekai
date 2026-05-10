package server

import (
	"io/fs"
	"net/http"
	"sekai/internal/handler"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Start(cfg *Config) error {
	secretKey, err := handler.GetJWTSecret()
	if err != nil {
		return err
	}
	hdlr := &handler.SessionHandler{
		JWTSecret: secretKey,
	}
	pgHdlr := &handler.PageHandler{
		JWTSecret: secretKey,
	}
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	publicFS, _ := fs.Sub(cfg.FS, "public")
	r.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.FS(publicFS))))

	// we will implement these handlers in the next sections
	r.Get("/", pgHdlr.HomePageHandler)
	r.Get("/about", pgHdlr.AboutPageHandler)
	r.Get("/pages/welcome", pgHdlr.WelcomePageHandler)
	r.Get("/pages/login", pgHdlr.LoginPageHandler)
	r.Post("/pages/login", pgHdlr.LoginFormHandler)

	r.Post("/sessions/login", hdlr.LoginHandler)
	r.Post("/sessions/refresh", hdlr.RefreshHandler)
	r.Post("/sessions/logout", hdlr.LogoutHandler)

	// start the server on port 8000
	return http.ListenAndServe(cfg.Addr, r)
}
