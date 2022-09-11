package main

import (
	"net/http"
	"github.com/jellyqwq/Paimon/cqtotg"
	
)

func main() {
	http.HandleFunc("/cq/", cqtotg.Post)
	http.ListenAndServe("127.0.0.1:6700", nil)
}