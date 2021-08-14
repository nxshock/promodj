package main

import (
	"fmt"
	"net/http"
)

func handleGenres(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	data := struct {
		Domain string
		Genres []Genre
	}{
		Domain: r.Host,
		Genres: Genres}

	err := templates.Lookup("genres.html").Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGetM3u(w http.ResponseWriter, r *http.Request) {
	genreCode := r.FormValue("genre")
	if genreCode == "" {
		http.Error(w, `"genre" field is not specified`, http.StatusBadRequest)
		return
	}

	tracks, err := tracksByGenre(genreCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := tracksToM3u(r.Host, tracks)

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.m3u8"`, genreCode))
	w.Header().Set("Content-Type", "audio/x-mpegurl")
	w.Header().Set("Accept-Ranges", "none")

	w.Write(b)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	stream(r.FormValue("url"), w)
}
