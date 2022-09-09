package main

import (
	"log"
	// news "github.com/jellyqwq/Paimon/news"
	"github.com/jellyqwq/Paimon/what"
)

func main() {
	s := what.BaiduBaike()

	log.Printf("%v", s)
}