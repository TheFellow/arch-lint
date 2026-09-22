package linter

import (
	"testing"

	"github.com/TheFellow/arch-lint/pkg/config"
)

func TestCheckImportForbidOrder(t *testing.T) {
	for _, exception := range []string{"except", "exempt"} {
		for _, reverse := range []bool{false, true} {
			for _, sameDomain := range []bool{false, true} {
				rules := config.Rules{Forbid: []string{"**", "app/domains/{domain}/**"}}
				if reverse {
					rules.Forbid[0], rules.Forbid[1] = rules.Forbid[1], rules.Forbid[0]
				}
				current := "app/domains/orders/service"
				if !sameDomain {
					current = "app/domains/payments/service"
				}
				if exception == "except" {
					rules.Except = []string{"app/domains/{domain}/**"}
				} else {
					rules.Exempt = []string{"app/domains/{domain}/**"}
				}
				got := CheckImport(config.Spec{Name: "boundaries", Rules: rules}, current, "app/domains/orders/model")
				wantViolation := exception == "except" && !sameDomain
				if (got != nil) != wantViolation {
					t.Errorf("%s reverse=%v sameDomain=%v: got %v, want violation=%v", exception, reverse, sameDomain, got, wantViolation)
				}
			}
		}
	}
}

func TestCheckImportPatternSemantics(t *testing.T) {
	tests := []struct {
		name              string
		rules             config.Rules
		current, imported string
		violation         bool
	}{
		{"stdlib catch all", config.Rules{Forbid: []string{"**"}}, "app", "fmt", true},
		{"stdlib exempt", config.Rules{Forbid: []string{"**"}, Exempt: []string{"fmt"}}, "app", "fmt", false},
		{"literal dot", config.Rules{Forbid: []string{"a.c/x"}}, "app", "abc/x", false},
		{"domain sibling", config.Rules{Forbid: []string{"app/domains/orders/**"}}, "app", "app/domains/orders-v2", false},
		{"forbid alternatives", config.Rules{Forbid: []string{"example/{beta,delta}/**"}}, "app", "example/delta/model", true},
		{"except alternatives", config.Rules{Forbid: []string{"**"}, Except: []string{"example/{beta,delta}/**"}}, "example/beta", "fmt", false},
		{"exempt alternatives", config.Rules{Forbid: []string{"**"}, Exempt: []string{"example/{beta,delta}/**"}}, "app", "example/delta", false},
		{"unbound positive", config.Rules{Forbid: []string{"**"}, Except: []string{"app/{domain}/**"}}, "app/orders", "fmt", true},
		{"unbound negative", config.Rules{Forbid: []string{"**"}, Except: []string{"app/{!domain}/**"}}, "app/orders", "fmt", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckImport(config.Spec{Name: "test", Rules: tt.rules}, tt.current, tt.imported)
			if (got != nil) != tt.violation {
				t.Fatalf("got %v, want violation=%v", got, tt.violation)
			}
		})
	}
}

func TestCheckImportNegatedExceptionOrder(t *testing.T) {
	for _, forbid := range [][]string{
		{"**", "app/{domain}/**"}, {"app/{domain}/**", "**"},
	} {
		spec := config.Spec{Rules: config.Rules{Forbid: forbid, Except: []string{"app/{!domain}/**"}}}
		for _, current := range []string{"app/orders", "app/payments"} {
			got := CheckImport(spec, current, "app/orders/model")
			if (got != nil) != (current == "app/orders") {
				t.Errorf("forbid=%v current=%s: got %v", forbid, current, got)
			}
		}
	}
}

func TestCheckImportDoesNotMergeCaptures(t *testing.T) {
	spec := config.Spec{Rules: config.Rules{
		Forbid: []string{"{first}/*", "*/{second}"},
		Except: []string{"{first}/{second}"},
	}}
	if got := CheckImport(spec, "orders/model", "orders/model"); got == nil {
		t.Fatal("exception must bind all variables from a single forbid match")
	}
}
