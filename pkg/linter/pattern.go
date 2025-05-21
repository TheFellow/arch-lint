package linter

import (
	"fmt"
	"regexp"
	"strings"
)

func matchPattern(pattern, path string) (map[string]string, bool) {
	// Convert gorilla-mux style pattern to regex
	regexPattern := "^" + regexp.QuoteMeta(pattern) + "$"
	regexPattern = strings.ReplaceAll(regexPattern, `\{`, `(?P<`)
	regexPattern = strings.ReplaceAll(regexPattern, `\}`, `>[^/]+)`)
	regexPattern = strings.ReplaceAll(regexPattern, `\*\*`, `.*`)
	regexPattern = strings.ReplaceAll(regexPattern, `\*`, `[^/]*`)

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, false
	}

	// Match the path against the regex
	match := re.FindStringSubmatch(path)
	if match == nil {
		return nil, false
	}

	// Extract named groups into a map
	vars := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && name != "" {
			vars[name] = match[i]
		}
	}
	return vars, true
}

func replaceVariables(pattern string, vars map[string]string) string {
	for key, value := range vars {
		placeholder := fmt.Sprintf("{%s}", key)
		pattern = strings.ReplaceAll(pattern, placeholder, value)
	}
	return pattern
}
