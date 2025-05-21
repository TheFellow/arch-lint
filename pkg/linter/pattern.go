package linter

import (
	"fmt"
	"regexp"
	"strings"
)

func matchPattern(pattern, path string) (map[string]string, bool) {
	// Split the pattern into segments
	segments := strings.Split(pattern, "/")
	for i, segment := range segments {
		// Handle variables
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			// Convert {var} to (?P<var>[^/]+)
			segment = fmt.Sprintf("(?P<%s>[^/]+)", segment[1:len(segment)-1])
		}

		// Handle single-level wildcards
		if segment == "*" {
			// Convert * to [^/]+
			segment = "[^/]+"
		}

		// Handle multi-level wildcards
		if segment == "**" {
			// Convert ** to .*
			segment = ".*"
		}

		// Update the segment
		segments[i] = segment
	}

	// Join the segments back together
	regexPattern := strings.Join(segments, "/")
	// Special case for /** at the end of the pattern
	if strings.HasSuffix(regexPattern, "/.*") {
		// If the pattern ends with a wildcard, allow empty string at the end
		regexPattern = strings.TrimSuffix(regexPattern, "/.*") + "/?.*"
	}
	regexPattern = "^" + regexPattern + "$"

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
		// Handle negated variables
		negatedPlaceholder := fmt.Sprintf("{!%s}", key)
		if strings.Contains(pattern, negatedPlaceholder) {
			// Replace negated variable with a regex that excludes the value
			pattern = strings.ReplaceAll(pattern, negatedPlaceholder,
				fmt.Sprintf("(?!%s)[^/]+", regexp.QuoteMeta(value)))
		} else {
			// Replace normal variables
			placeholder := fmt.Sprintf("{%s}", key)
			pattern = strings.ReplaceAll(pattern, placeholder, value)
		}
	}
	return pattern
}
