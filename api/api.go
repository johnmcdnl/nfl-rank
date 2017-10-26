package main

import (
	"flag"
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/johnmcdnl/nfl-rank/api/resources"
)

func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			// We set a default error status response code if one hasn't been set.
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			// We log the error
			fmt.Printf("Logging err: %s\n", err.Error())

			// We change the response to not reveal the actual error message,
			// instead we can transform the message something more friendly or mapped
			// to some code / language, etc.
			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
}

func main() {
	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/articles", func(r chi.Router) {
		r.Get("/", resources.ListArticles)
		r.Route("/{articleID}", func(r chi.Router) {
			r.Get("/", resources.GetArticle)
		})
	})

	r.Route("/team-ratings", func(r chi.Router) {
		r.Get("/", resources.ListArticles)
		r.Route("/{teamID}", func(r chi.Router) {
			r.Get("/", resources.GetTeamRating)
		})
	})

	http.ListenAndServe(":3333", r)
}
