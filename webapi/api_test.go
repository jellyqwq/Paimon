package webapi

import (
	// "encoding/json"
	// "bytes"
	// "fmt"
	// "time"
	// "log"
	// "strconv"

	// "fmt"
	// "log"
	// "io"
	// "net/http"
	// "net/url"

	// "regexp"
	// "strings"
	"fmt"
	"testing"

	// "github.com/jellyqwq/Paimon/webapi"
	// "github.com/jellyqwq/Paimon/requests"
	// "github.com/jellyqwq/Paimon/tools"
)

func Test(t *testing.T) {
	r, e := RranslateByYouDao("Never imagine how many bugs you need to fix.")
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(r)
}