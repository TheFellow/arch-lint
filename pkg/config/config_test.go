package config

import (
	"os"
	"strings"
	"testing"

	"github.com/TheFellow/arch-lint/pkg/testutil"
)

func TestLoad_SchemaValidation(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	valid := []byte("specs:\n  - name: test\n    packages:\n      include:\n        - pkg\n    rules:\n      forbid:\n        - other\n")
	good := dir + "/good.yml"
	os.WriteFile(good, valid, 0o644)

	_, err := Load(good)
	testutil.Equals(t, err, nil)

	bad := dir + "/bad.yml"
	os.WriteFile(bad, []byte("foo: bar"), 0o644)
	_, err = Load(bad)
	testutil.ErrorIf(t, err == nil, "expected error")
}

func TestLoadPatternValidation(t *testing.T) {
	tests := []struct {
		field, pattern string
		valid          bool
	}{
		{"forbid", "example/{beta,delta}/**", true},
		{"except", "example/{!domain}/**", true},
		{"exempt", "example/{beta,delta}/**", true},
		{"forbid", "example/{beta,}/**", false},
		{"forbid", "example/{bad-name}/**", false},
		{"forbid", "example/{!domain}/**", false},
		{"except", "example/{broken/**", false},
		{"exempt", "example/{domain}/{domain}", false},
		{"forbid", "example/prefix*", false},
		{"include", "example/{beta,delta}/**", true},
		{"exclude", "example/[", false},
	}
	for _, tt := range tests {
		t.Run(tt.field+":"+tt.pattern, func(t *testing.T) {
			cfg := "specs:\n  - name: boundaries\n    packages:\n      include: [example/**]\n"
			if tt.field == "include" {
				cfg = "specs:\n  - name: boundaries\n    packages:\n      include: [\"" + tt.pattern + "\"]\n"
			}
			if tt.field == "exclude" {
				cfg += "      exclude: [\"" + tt.pattern + "\"]\n"
			}
			cfg += "    rules:\n"
			if tt.field != "forbid" {
				cfg += "      forbid: [\"**\"]\n"
			}
			if tt.field != "include" && tt.field != "exclude" {
				cfg += "      " + tt.field + ": [\"" + tt.pattern + "\"]\n"
			}
			path := t.TempDir() + "/rules.yml"
			if err := os.WriteFile(path, []byte(cfg), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := Load(path)
			if (err == nil) != tt.valid {
				t.Fatalf("Load: %v, want valid=%v", err, tt.valid)
			}
			if !tt.valid && (!strings.Contains(err.Error(), "boundaries") || !strings.Contains(err.Error(), tt.field) || !strings.Contains(err.Error(), tt.pattern)) {
				t.Fatalf("error lacks pattern context: %v", err)
			}
		})
	}
}
