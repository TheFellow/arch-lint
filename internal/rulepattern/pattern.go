// Package rulepattern compiles import-rule patterns. Package selection uses
// doublestar instead, because braces in rule patterns also represent captures.
package rulepattern

import (
	"fmt"
	"regexp"
	"strings"
)

var variableName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Regex translates a rule pattern to an anchored regular expression.
// allowNegation permits {!variable} references in except and exempt rules.
func Regex(pattern string, allowNegation bool) (string, error) {
	segments := make([]string, 0)
	for _, segment := range strings.Split(pattern, "/") {
		if segment == "**" && len(segments) > 0 && segments[len(segments)-1] == "**" {
			continue
		}
		segments = append(segments, segment)
	}
	var result strings.Builder
	result.WriteString("^")
	seen := make(map[string]bool)
	for i, segment := range segments {
		if segment == "**" {
			if i == len(segments)-1 {
				if i == 0 {
					result.WriteString(".*")
				} else {
					result.WriteString("(?:/.*)?")
				}
			} else {
				if i > 0 {
					result.WriteString("/")
				}
				result.WriteString("(?:[^/]+/)*")
			}
			continue
		}
		if i > 0 && segments[i-1] != "**" {
			result.WriteString("/")
		}
		switch {
		case segment == "*":
			result.WriteString("[^/]+")
		case strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}"):
			body := segment[1 : len(segment)-1]
			if strings.Contains(body, ",") {
				alternatives := strings.Split(body, ",")
				for j, alternative := range alternatives {
					if alternative == "" || strings.ContainsAny(alternative, "{}*!") {
						return "", fmt.Errorf("invalid alternative %q", alternative)
					}
					alternatives[j] = regexp.QuoteMeta(alternative)
				}
				result.WriteString("(?:" + strings.Join(alternatives, "|") + ")")
			} else {
				if allowNegation {
					body = strings.TrimPrefix(body, "!")
				}
				if !variableName.MatchString(body) {
					return "", fmt.Errorf("invalid variable name %q", body)
				}
				if seen[body] {
					return "", fmt.Errorf("duplicate variable %q", body)
				}
				seen[body] = true
				result.WriteString("(?P<" + body + ">[^/]+)")
			}
		default:
			if strings.ContainsAny(segment, "{}*") {
				return "", fmt.Errorf("wildcards and captures must occupy a complete segment: %q", segment)
			}
			result.WriteString(regexp.QuoteMeta(segment))
		}
	}
	result.WriteString("$")
	return result.String(), nil
}
