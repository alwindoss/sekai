package main

import (
	"embed"
	"io/fs"
	"net/http"
	"sekai/data"
	"sekai/templates/components"
	"sekai/templates/layouts"
	"sekai/templates/pages"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed public/*
var embedFS embed.FS

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", templ.Handler(layouts.Default(pages.Index(), &data.PageData{
		Title: "Home",
	})).ServeHTTP)
	r.Get("/about", templ.Handler(layouts.Default(pages.About(), &data.PageData{
		Title: "About",
	})).ServeHTTP)
	r.Post("/clicked", templ.Handler(components.TimeNow()).ServeHTTP)
	publicFS, _ := fs.Sub(embedFS, "public")
	r.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.FS(publicFS))))
	http.ListenAndServe(":3000", r)
}
