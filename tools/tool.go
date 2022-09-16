package tools

import (
	"regexp"
	"crypto/md5"
	"fmt"
)

func GetParamsOneDimension(compile *regexp.Regexp, str string) (paramsMap map[string]string){
	match := compile.FindStringSubmatch(str)

	paramsMap = make(map[string]string)
	for i, name := range compile.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func IsOneDimensionSliceContainsString(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
