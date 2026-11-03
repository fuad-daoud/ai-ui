package ui_test

import (
	"os"
	"regexp"
	"testing"
)

// documentedClasses is every class DESIGN.md documents for v0.1. Renaming
// or removing one is a breaking change (spec §9).
var documentedClasses = []string{
	"ui-btn",
	"ui-field", "ui-field-label", "ui-field-hint", "ui-field-error",
	"ui-input", "ui-select", "ui-textarea", "ui-form-row",
	"ui-check", "ui-check-input",
	"ui-pill",
	"ui-dialog", "ui-dialog-head", "ui-dialog-title", "ui-dialog-x", "ui-dialog-body", "ui-dialog-actions",
	"ui-callout", "ui-callout-title",
	"ui-table-wrap", "ui-table",
	"ui-disclosure", "ui-disclosure-summary", "ui-disclosure-body",
	"ui-card",
	"ui-muted", "ui-mono", "ui-num", "ui-vh", "ui-stack", "ui-cluster",
}

func TestDocumentedClassesExist(t *testing.T) {
	css, err := os.ReadFile("dist/ui.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range documentedClasses {
		if !regexp.MustCompile(`\.` + regexp.QuoteMeta(c) + `[\s,.:{\[>)]`).Match(css) {
			t.Errorf(".%s is documented but not defined in dist/ui.css", c)
		}
	}
}
