package express

import (
	"strings"
)

func findWildcards(path string) []string {
	params := []string{}
	segments := strings.Split(path, "/")
	for _, s := range segments {
		if s == "" {
			continue
		}
		if strings.Contains(s, ":") {
			s = strings.Replace(s, ":", "", -1)
			s = strings.Replace(s, "?", "", -1)
			params = append(params, s)
		} else if strings.Contains(s, "*") {
			params = append(params, "*")
		}
	}
	return params
}

func pathToRegex(path string) (regex string) {
	regex += "^"
	segments := strings.Split(path, "/")
	for _, s := range segments {
		if s == "" {
			continue
		}
		if strings.Contains(s, ":") && strings.Contains(s, "?") {
			regex += "(?:/([^/]+?))?"
		} else if strings.Contains(s, ":") {
			regex += "/(?:([^/]+?))"
		} else if strings.Contains(s, "*") {
			regex += "/(.*)"
		} else {
			regex += "/" + s
		}
	}
	regex += "/?$"
	return regex
}
