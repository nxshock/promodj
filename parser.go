package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Genre represents track information
type TrackInfo struct {
	Title string
	Url   string
}

// Genre represents genre information
type Genre struct {
	Name string
	Code string
}

// Genres holds cached list of available genres
var Genres []Genre

func UpdateGenres() error {
	var err error

	Genres, err = updateGenreList()
	if err != nil {
		return fmt.Errorf("get genres list failed: %w", err)
	}

	return nil
}

func updateGenreList() ([]Genre, error) {
	url := "https://promodj.com/music"

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var genres []Genre
	doc.Find("div.styles_tagcloud > a").Each(func(i int, s *goquery.Selection) {
		genres = append(genres, Genre{s.Text(), strings.TrimPrefix(s.AttrOr("href", ""), "/music/")})
	})

	return genres, nil
}

// tracksByGenre возвращает список треков по указанному жанру
func tracksByGenre(genre string, params url.Values) ([]TrackInfo, error) {
	if params == nil {
		params = url.Values{"download": []string{"1"}}
	} else {
		params.Set("download", "1") // only available tracks
	}

	var result []TrackInfo

	for i := 1; i <= 50; i++ {
		params.Set("page", strconv.Itoa(i))
		//url := fmt.Sprintf("https://promodj.com/music/%s?download=1&page=%d", genre, i)
		url := constructUrl(genre, params)

		doc, err := goquery.NewDocument(url)
		if err != nil {
			return nil, err
		}

		doc.Find("div.title > a.invert").Each(
			func(n int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if !exists {
					return
				}

				fields := strings.Split(href, "/")
				if len(fields) != 7 {
					return
				}

				fileUrl := fmt.Sprintf("https://promodj.com/download/%s/%s.mp3", fields[5], fields[6])

				result = append(result, TrackInfo{s.Text(), fileUrl})
			})
	}

	result = removeDuplicate(result)

	return result, nil
}

// tracksToM3u возвращает байты M3U-плейлиста, сгенерированного по указанному
// списку треков
func tracksToM3u(host string, tracks []TrackInfo) []byte {
	b := new(bytes.Buffer)

	b.Write([]byte{0xEF, 0xBB, 0xBF})
	fmt.Fprint(b, "#EXTM3U\n")

	for _, track := range tracks {
		fmt.Fprintf(b, "#EXTINF:-1,%s\n", track.Title)

		u, _ := url.Parse(fmt.Sprintf("https://%s/stream", host))
		q := make(url.Values)
		q.Add("url", track.Url)
		u.RawQuery = q.Encode()
		fmt.Fprintf(b, "%s\n", u.String())
	}

	return b.Bytes()
}

func removeDuplicate(strSlice []TrackInfo) []TrackInfo {
	allKeys := make(map[string]bool)
	list := []TrackInfo{}
	for _, item := range strSlice {
		if _, value := allKeys[item.Url]; !value {
			allKeys[item.Url] = true
			list = append(list, item)
		}
	}
	return list
}

func constructUrl(genre string, params url.Values) string {
	urlTemplate := fmt.Sprintf("https://promodj.com/music/%s", genre)

	if genre == "" {
		urlTemplate = "https://promodj.com/music"
	}

	u, err := url.Parse(urlTemplate)
	if err != nil {
		panic(err)
	}

	u.RawQuery = params.Encode()

	//?download=1&page=%d
	return u.String()
}
