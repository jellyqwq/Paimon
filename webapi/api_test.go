package webapi

import (
	// "encoding/json"
	// "bytes"
	// "fmt"
	"log"
	"strconv"

	// "log"
	// "io"
	// "net/http"
	// "net/url"

	"regexp"
	"testing"

	// "github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/tools"
)

func Test(t *testing.T) {
	str := `3 minutes, 51 seconds`
	compileHMS := regexp.MustCompile(`(?:(?P<H>[0-9]+) hours?, )?(?:(?P<M>[0-9]+) minutes?, )?(?P<S>[0-9]+) seconds?`)
	paramsMap := tools.GetParamsOneDimension(compileHMS, str)

	seconds, err := strconv.Atoi(paramsMap["S"])
	if err != nil {
		log.Println(err)
		return
	}

	if paramsMap["H"] != "" {
		h, err := strconv.Atoi(paramsMap["H"])
		if err != nil {
			log.Println(err)
			return
		}
		seconds += 3600 * h
	}

	if paramsMap["M"] != "" {
		m, err := strconv.Atoi(paramsMap["M"])
		if err != nil {
			log.Println(err)
			return
		}
		seconds += 60 * m
	}

	log.Println(seconds)
}