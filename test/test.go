package main

import (
	"log"
	news "github.com/jellyqwq/Paimon/news"
)

func main() {
	s, err := news.WeiboHotWords()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v", s)
}