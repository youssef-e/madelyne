package unittester

import (
	"regexp"
	"strings"
)

var replaceVarRegexp = regexp.MustCompile(`\#(.*?)\#`)

func ReplaceWithEnvValue(src string, env map[string]string) string {
	found := replaceVarRegexp.FindAllStringSubmatch(src, -1)
	for _, f := range found {
		for _, fvalue := range f {
			v, ok := env[fvalue]
			if ok {
				src = strings.ReplaceAll(src, "#"+fvalue+"#", v)
			}
		}
	}
	return src
}
