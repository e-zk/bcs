package bcs

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Track struct {
	Number int // track number
	Title  string
	Artist string
	Time   string // duration (not impl)
	Url    string // url
}

type Album struct {
	URL           string
	BaseURL       string
	Title         string
	Artist        string
	Tracks        []Track
	CoverLargeURL string
	CoverURL      string
	AlbumAbout    string
	AlbumCredits  string
	Tags          []string
}

func ParseAlbum(albumURL string) (album Album, err error) {
	album.URL = albumURL

	// parse url to get base
	u, err := url.Parse(albumURL)
	if err != nil {
		return
	}
	album.BaseURL = u.Scheme + "://" + u.Host

	r, err := http.Get(album.URL)
	if err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return album, fmt.Errorf("status code error: %d %s", r.StatusCode, r.Status)
	}

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return
	}

	album.Title = strings.TrimSpace(doc.Find(".trackView .trackTitle").Text())
	album.Artist = strings.TrimSpace(doc.Find(".trackView h3 > span > a").Text())
	var tracks []Track
	doc.Find(".track_list .track_row_view").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".title span.track-title").Text()
		titleUrl := s.Find(".title > a").AttrOr("href", "")
		num, err := strconv.Atoi(strings.Split(s.Find(".track_number").Text(), ".")[0])
		if err != nil {
			err = fmt.Errorf("error parsing track number: %v", err)
			return
		}
		t := Track{
			Url:    titleUrl,
			Number: num,
			Title:  title,
			Artist: album.Artist,
		}
		tracks = append(tracks, t)
	})
	if err != nil {
		return
	}

	album.Tracks = tracks
	album.CoverURL = doc.Find("#tralbumArt > a.popupImage > img").AttrOr("src", "")
	album.CoverLargeURL = doc.Find("#tralbumArt > a.popupImage").AttrOr("href", "")
	album.AlbumAbout = strings.TrimSpace(doc.Find(".tralbum-about").Text())
	album.AlbumCredits = strings.TrimSpace(doc.Find(".tralbum-credits").Text())
	return
}
