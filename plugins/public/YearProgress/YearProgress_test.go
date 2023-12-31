package YearProgress

import (
	"fmt"
	"testing"
)

// var maxDays int64
func Test(t *testing.T) {
	ypc := NewYearProgressConfig()
	fmt.Println(ypc.display())
	fmt.Println(ypc)
	fmt.Println(ypc.GetYearProgress())
	fmt.Println(ypc)
	fmt.Println(ypc.GetYearProgress())
}


