package webapi

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	fmt.Println(int64(time.Now().Local().Month()))
}
