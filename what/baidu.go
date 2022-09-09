package what

import (
	"fmt"
	// "net/http"
	"regexp"
)

func BaiduBaike() (string) {
	a := "Ely你好"
	reg := regexp.MustCompile(`^(爱|e|E){1}(莉|ly){1}(希雅|sia)?`)
	return fmt.Sprintf("%s\n", reg.ReplaceAll([]byte(a), []byte("")))
}