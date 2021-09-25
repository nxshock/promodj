package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var templates *template.Template

func init() {
	log.SetFlags(0)

	if len(os.Args) == 2 {
		if err := initConfig(os.Args[1]); err != nil {
			log.Fatalln("config error:", err)
		}
	} else {
		if err := initConfig(defaultConfigFilePath); err != nil {
			log.Fatalln("config error:", err)
		}
	}

	err := initTepmplates()
	if err != nil {
		log.Fatalln(err)
	}

	err = UpdateGenres()
	if err != nil {
		log.Fatalln(err)
	}

	http.DefaultClient.Timeout = 5 * time.Second
}

func initTepmplates() error {
	var err error

	templates, err = template.ParseFS(siteFS, "site/*.htm")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	http.HandleFunc("/", handleGenres)
	http.HandleFunc("/genres", handleGenres)
	http.HandleFunc("/getm3u", handleGetM3u)
	http.HandleFunc("/stream", handleStream)

	err := http.ListenAndServe(config.ListenAddr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
