package regex

import (
	"regexp"
)

func EvalRegex(r *regexp.Regexp, s string) map[string]string {
	matches := r.FindAllStringSubmatch(s, -1)
  retMap := make(map[string]string)
	for _, m := range matches {
		for i, v := range m {
			if i != 0 && v != ""{
				retMap[r.SubexpNames()[i]] = v
			}
		}
	}
  return retMap
}
