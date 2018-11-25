package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"html/template"
	"net/http"
	"video-stream/controller"
)

type VideoPage struct {
	Filename string
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Route("/web", func(r chi.Router) {
		r.Get("/*", func (w http.ResponseWriter, r *http.Request) {
			videoPath := chi.URLParam(r, "*")
			videoPath = "http://" + r.Host + "/api/" + videoPath
			videoPage := VideoPage{videoPath}
			t, _ := template.ParseFiles("static/video.html")
			t.Execute(w, videoPage)
		})
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/*", controller.VideoStream)
	})

	http.ListenAndServe(":8000", r)
}