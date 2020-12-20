package unittester

import (
	"bytes"
	"regexp"
	"strings"
)

var replaceVarRegexp = regexp.MustCompile(`\#(.*?)\#`)

func ReplaceStringWithEnvValue(src string, env map[string]string) string {
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

func ReplaceWithEnvValue(src []byte, env map[string]string) []byte {
	found := replaceVarRegexp.FindAllSubmatch(src, -1)
	for _, f := range found {
		for _, fvalue := range f {
			v, ok := env[string(fvalue)]
			if ok {
				search := []byte("#")
				search = append(search, fvalue...)
				search = append(search, []byte("#")...)
				src = bytes.ReplaceAll(src, search, []byte(v))
			}
		}
	}
	return src
}
