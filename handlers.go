package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func handleGenres(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.FileServer(http.FS(stripSiteFS)).ServeHTTP(w, r)
		return
	}

	data := struct {
		Domain string
		Genres []Genre
	}{
		Domain: r.Host,
		Genres: Genres}

	err := templates.Lookup("index.htm").Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePlayer(w http.ResponseWriter, r *http.Request) {
	genreCode := r.FormValue("genre")
	params := url.Values{}

	if r.FormValue("top100") != "" {
		params.Set("top100", "1")
	}

	tracks, err := tracksByGenre(genreCode, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type J struct {
		Title string
		File  string
	}

	var data []J

	for _, track := range tracks {
		u, _ := url.Parse("stream")
		q := make(url.Values)
		q.Add("url", track.Url)
		u.RawQuery = q.Encode()

		data = append(data, J{track.Title, u.String()})
	}

	err = templates.Lookup("player.htm").Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGetM3u(w http.ResponseWriter, r *http.Request) {
	genreCode := r.FormValue("genre")
	params := url.Values{}

	if r.FormValue("top100") != "" {
		params.Set("top100", "1")
	}

	tracks, err := tracksByGenre(genreCode, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := tracksToM3u(r.Host, tracks)

	if genreCode == "" {
		genreCode = "music"
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.m3u8"`, genreCode))
	w.Header().Set("Content-Type", "audio/x-mpegurl")
	w.Header().Set("Accept-Ranges", "none")

	_, _ = w.Write(b)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	bitrateKbps, err := strconv.ParseUint(r.URL.Query().Get("b"), 10, 64)
	if err != nil {
		bitrateKbps = config.Bitrate
	}

	err = stream(r.FormValue("url"), w, bitrateKbps)
	if err != nil {
		log.Println(err)
	}
}
