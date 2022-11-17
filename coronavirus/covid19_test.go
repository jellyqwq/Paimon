package coronavirus

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	a, e := MainHandle()
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(a.ProvinceInlineKeyborad[0])
	}
}