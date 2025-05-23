package linter

import (
	"testing"

	"github.com/TheFellow/arch-lint/pkg/testutil"
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
			pattern:   "example/beta/bookstore/app/{feature}/some/path",
			path:      "example/beta/bookstore/app/feature1/some/path",
			wantVars:  map[string]string{"feature": "feature1"},
			wantMatch: true,
		},
		{
			name:      "match with multiple variables",
			pattern:   "example/{section}/bookstore/{feature}/some/path",
			path:      "example/beta/bookstore/feature1/some/path",
			wantVars:  map[string]string{"section": "beta", "feature": "feature1"},
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
			name:      "match with trailing multi-level wildcard",
			pattern:   "example/beta/bookstore/**",
			path:      "example/beta/bookstore/feature1/some/path",
			wantVars:  nil,
			wantMatch: true,
		},
		{
			name:      "match with multi-level wildcard in the middle",
			pattern:   "example/beta/bookstore/**/api",
			path:      "example/beta/bookstore/feature1/some/path/api",
			wantVars:  nil,
			wantMatch: true,
		},
		{
			name:      "no match with multi-level wildcard in the middle",
			pattern:   "example/beta/bookstore/**/api",
			path:      "example/beta/bookstore/feature1/some/path/util",
			wantVars:  nil,
			wantMatch: false,
		},
		{
			name:      "match with multi-level wildcard at the end does not need to match anything",
			pattern:   "example/beta/bookstore/**",
			path:      "example/beta/bookstore",
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

func TestExceptRegex(t *testing.T) {
	tests := []struct {
		name      string
		pattern   string
		vars      map[string]string
		path      string
		wantMatch bool
	}{
		{
			name:      "simple match with variable",
			pattern:   "example/beta/bookstore/app/{feature}/some/path",
			vars:      map[string]string{"feature": "feature1"},
			path:      "example/beta/bookstore/app/feature1/some/path",
			wantMatch: true,
		},
		{
			name:      "no match with different path",
			pattern:   "example/beta/bookstore/app/{feature}/some/path",
			vars:      map[string]string{"feature": "feature1"},
			path:      "example/beta/bookstore/app/feature2/some/path",
			wantMatch: false,
		},
		{
			name:      "match with negated variable",
			pattern:   "example/beta/bookstore/app/{!feature}/some/path",
			vars:      map[string]string{"feature": "feature1"},
			path:      "example/beta/bookstore/app/feature2/some/path",
			wantMatch: true,
		},
		{
			name:      "no match with negated variable",
			pattern:   "example/beta/bookstore/app/{!feature}/some/path",
			vars:      map[string]string{"feature": "feature1"},
			path:      "example/beta/bookstore/app/feature1/some/path",
			wantMatch: false,
		},
		{
			name:      "match with multi-level wildcard",
			pattern:   "example/beta/bookstore/**",
			vars:      nil,
			path:      "example/beta/bookstore/feature1/some/path",
			wantMatch: true,
		},
		{
			name:      "no match with empty path",
			pattern:   "example/beta/bookstore/**",
			vars:      nil,
			path:      "",
			wantMatch: false,
		},
		{
			name:      "match with single-level wildcard",
			pattern:   "example/beta/bookstore/*/some/path",
			vars:      nil,
			path:      "example/beta/bookstore/feature1/some/path",
			wantMatch: true,
		},
		{
			name:      "no match with single-level wildcard",
			pattern:   "example/beta/bookstore/*/some/path",
			vars:      nil,
			path:      "example/beta/bookstore/feature1/other/path",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch := exceptRegex(tt.pattern, tt.path, tt.vars)
			testutil.Equals(t, gotMatch, tt.wantMatch)
		})
	}
}
