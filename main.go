package main

import (
	"net/http"
	"sekai/data"
	"sekai/templates/layouts"
	"sekai/templates/pages"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", templ.Handler(layouts.Default(pages.Index(), &data.PageData{
		Title: "Home",
	})).ServeHTTP)
	r.Get("/about", templ.Handler(layouts.Default(pages.About(), &data.PageData{
		Title: "About",
	})).ServeHTTP)
	http.ListenAndServe(":3000", r)
}
