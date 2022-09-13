package tools

import (
	"regexp"
)

func GetParamsOneDimension(compile *regexp.Regexp, s string) (paramsMap map[string]string){
	match := compile.FindStringSubmatch(s)

	paramsMap = make(map[string]string)
	for i, name := range compile.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}