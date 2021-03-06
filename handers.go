package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/nxshock/promodj/api"
)

var templates *template.Template

func init() {
	var err error

	templateFiles, err := filepath.Glob(filepath.Join(templateFilesPath, "*.html"))
	if err != nil {
		log.Fatalln(err)
	}

	templates, err = template.ParseFiles(templateFiles...)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/playlist", handlePlaylist)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/genres", handleGenres)
	http.HandleFunc("/getRandomTrackInfoByGenre", handleGetRandomTrackInfoByGenre)
	http.HandleFunc("/getRandomTrackDataByGenre", handleGetRandomTrackDataByGenre)
	http.Handle("/", http.FileServer(http.Dir(siteFilesPath)))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	stream(r.FormValue("url"), w)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	download(r.FormValue("url"), "file.opus", w)
}

func handlePlaylist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", `attachment; filename="playlist.m3u"`)
	w.Header().Set("Content-Type", config.ContentType)
	w.Header().Set("Accept-Ranges", "none")

	host := r.Host
	if _, port, err := net.SplitHostPort(config.ListenAddr); err != nil && port != "0" && port != "80" {
		host = host + ":" + port
	}

	w.Write(api.GetM3uPlaylist("https://" + host + "/getRandomTrackDataByGenre?g="))
}

func handleGetRandomTrackDataByGenre(w http.ResponseWriter, r *http.Request) {
	genre := r.FormValue("g")
	if genre == "" {
		http.Error(w, "genre field (g) can't be empty", http.StatusBadRequest)
		return
	}

	trackInfo, err := api.RandomTrackInfoByGenre(genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stream(trackInfo.DownloadUrl, w)
}

func handleGetRandomTrackInfoByGenre(w http.ResponseWriter, r *http.Request) {
	genre := r.FormValue("g")
	if genre == "" {
		http.Error(w, "genre field (g) can't be empty", http.StatusBadRequest)
		return
	}

	track, err := api.RandomTrackInfoByGenre(genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(b)
}

func handleGenres(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Domain string
		Genres []api.Genre
	}{
		Domain: r.Host,
		Genres: api.Genres}

	err := templates.Lookup("genres.html").Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	res, err := api.Search(r.FormValue("q"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = struct {
		Query   string
		Results []*api.TrackInfo
	}{r.FormValue("q"), res}

	templates.Lookup("search.html").Execute(w, data)
}
