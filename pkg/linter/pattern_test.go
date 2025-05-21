package linter

import (
	"testing"

	"github.com/TheFellow/go-arch-lint/pkg/testutil"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name      string
		pattern   string
		path      string
		wantVars  map[string]string
		wantMatch bool
	}{
		{
			name:      "simple match with variable",
			pattern:   "example/beta/bookstore/app/{feature}/**",
			path:      "example/beta/bookstore/app/feature1/some/path",
			wantVars:  map[string]string{"feature": "feature1"},
			wantMatch: true,
		},
		{
			name:      "no match due to different path",
			pattern:   "example/beta/bookstore/app/{feature}/**",
			path:      "example/beta/bookstore/api/feature1/some/path",
			wantVars:  nil,
			wantMatch: false,
		},
		{
			name:      "match with multiple variables",
			pattern:   "example/{section}/bookstore/{feature}/**",
			path:      "example/beta/bookstore/feature1/some/path",
			wantVars:  map[string]string{"section": "beta", "feature": "feature1"},
			wantMatch: true,
		},
		{
			name:      "match with single-level wildcard",
			pattern:   "example/beta/bookstore/*/some/path",
			path:      "example/beta/bookstore/feature1/some/path",
			wantVars:  nil,
			wantMatch: true,
		},
		{
			name:      "no match with single-level wildcard",
			pattern:   "example/beta/bookstore/*/some/path",
			path:      "example/beta/bookstore/feature1/other/path",
			wantVars:  nil,
			wantMatch: false,
		},
		{
			name:      "match with multi-level wildcard",
			pattern:   "example/beta/bookstore/**",
			path:      "example/beta/bookstore/feature1/some/path",
			wantVars:  nil,
			wantMatch: true,
		},
		{
			name:      "no match with empty path",
			pattern:   "example/beta/bookstore/**",
			path:      "",
			wantVars:  nil,
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVars, gotMatch := matchPattern(tt.pattern, tt.path)
			testutil.Equals(t, gotMatch, tt.wantMatch)
			testutil.Equals(t, gotVars, tt.wantVars, cmpopts.EquateEmpty())
		})
	}
}

func TestReplaceVariables(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		vars     map[string]string
		expected string
	}{
		{
			name:     "single variable replacement",
			pattern:  "example/{feature}/path",
			vars:     map[string]string{"feature": "beta"},
			expected: "example/beta/path",
		},
		{
			name:     "multiple variable replacements",
			pattern:  "example/{section}/{feature}/path",
			vars:     map[string]string{"section": "alpha", "feature": "gamma"},
			expected: "example/alpha/gamma/path",
		},
		{
			name:     "no variables to replace",
			pattern:  "example/static/path",
			vars:     map[string]string{},
			expected: "example/static/path",
		},
		{
			name:     "variable not in map",
			pattern:  "example/{missing}/path",
			vars:     map[string]string{"feature": "beta"},
			expected: "example/{missing}/path",
		},
		{
			name:     "empty pattern",
			pattern:  "",
			vars:     map[string]string{"feature": "beta"},
			expected: "",
		},
		{
			name:     "negated variable replacement",
			pattern:  "example/{!feature}/path",
			vars:     map[string]string{"feature": "beta"},
			expected: "example/(?!beta)[^/]+/path",
		},
		{
			name:     "negated variable not in map",
			pattern:  "example/{!missing}/path",
			vars:     map[string]string{"feature": "beta"},
			expected: "example/{!missing}/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceVariables(tt.pattern, tt.vars)
			testutil.Equals(t, result, tt.expected)
		})
	}
}
