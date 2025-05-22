package linter

import (
	"fmt"
	"regexp"
	"strings"
)

func matchPattern(pattern, path string) (map[string]string, bool) {
	regexPattern := escapePattern(pattern)

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

func escapePattern(pattern string) string {
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
	return regexPattern
}

func replaceVariables(pattern string, vars map[string]string) string {
	segments := strings.Split(pattern, "/")
	for i, segment := range segments {
		for key, value := range vars {
			// Handle negated variables
			negatedPlaceholder := fmt.Sprintf("{!%s}", key)
			if segment == negatedPlaceholder {
				var sb strings.Builder
				for _, c := range value {
					sb.WriteRune('[')
					sb.WriteRune('^')
					sb.WriteRune(c)
					sb.WriteRune(']')
				}
				segment = sb.String()
			}
			if segment == fmt.Sprintf("{%s}", key) {
				// Replace normal variables
				placeholder := fmt.Sprintf("{%s}", key)
				segment = strings.ReplaceAll(segment, placeholder, value)
			}
		}
		// Update the segment
		segments[i] = segment
	}

	// Recombine the segments into the final pattern
	return strings.Join(segments, "/")
}
