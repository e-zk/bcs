package main

import (
	"fmt"
	"log"
	"os"

	"go.zakaria.org/bcs"
)

// usage: ./bcs <album url from bandcamp>
func main() {
	if len(os.Args) < 2 {
		log.Fatal("no album given.")
	}
	alb, err := bcs.ParseAlbum(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("title:\t%s\n", alb.Title)
	fmt.Printf("artist:\t%s\n", alb.Artist)
	fmt.Printf("cover:\t%s\n", alb.CoverURL)
	fmt.Printf("cover:\t%s\n", alb.CoverLargeURL)
	fmt.Printf("tracks:\n")
	for _, t := range alb.Tracks {
		fmt.Printf("[%d] %s (%s%s)\n", t.Number, t.Title, alb.BaseURL, t.Url)
	}
	fmt.Printf("about:\n%s\n---\n", alb.AlbumAbout)
	fmt.Printf("credits:\n%s\n---\n", alb.AlbumCredits)
}
