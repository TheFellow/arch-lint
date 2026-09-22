package linter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/TheFellow/arch-lint/internal/rulepattern"
)

func MatchPattern(pattern, path string) (map[string]string, bool) {
	regexPattern := EscapePattern(pattern)

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

func EscapePattern(pattern string) string {
	regexPattern, err := rulepattern.Regex(pattern, false)
	if err != nil {
		// Keep the exported helper's signature; invalid patterns cannot match.
		return "(?!)"
	}
	return regexPattern
}

func ReplaceVariables(pattern string, vars map[string]string) string {
	segments := strings.Split(pattern, "/")
	for i, segment := range segments {
		for key := range vars {
			// Handle negated variables by treating them as normal variables in the regex
			negatedPlaceholder := fmt.Sprintf("{!%s}", key)
			if segment == negatedPlaceholder {
				segment = fmt.Sprintf("{%s}", key)
			}
		}
		segments[i] = segment
	}
	return strings.Join(segments, "/")
}

func ExceptRegex(pattern, path string, vars map[string]string) bool {
	// Every reference must be bound by this particular forbid match.
	for _, segment := range strings.Split(pattern, "/") {
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") && !strings.Contains(segment, ",") {
			name := strings.TrimPrefix(segment[1:len(segment)-1], "!")
			if _, ok := vars[name]; !ok {
				return false
			}
		}
	}
	// Replace variables in the pattern
	regexPattern := EscapePattern(ReplaceVariables(pattern, vars))

	// Compile the regex
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return false
	}

	// Match the path against the regex
	match := re.FindStringSubmatch(path)
	if match == nil {
		return false
	}

	// Extract named groups into a map
	capturedVars := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && name != "" {
			capturedVars[name] = match[i]
		}
	}

	// Validate variables (both positive and negated)
	for key, value := range vars {
		negatedPlaceholder := fmt.Sprintf("{!%s}", key)
		positivePlaceholder := fmt.Sprintf("{%s}", key)

		if strings.Contains(pattern, negatedPlaceholder) {
			// For negated variables, ensure the captured value does not match the forbidden value
			if capturedVars[key] == value {
				return false // Negated variable matches the forbidden value
			}
		} else if strings.Contains(pattern, positivePlaceholder) {
			// For positive variables, ensure the captured value matches the expected value
			if capturedVars[key] != value {
				return false // Positive variable does not match the expected value
			}
		}
	}

	return true
}
