package tests

import (
	"github.com/wyw14/cry052/internal/service"
	"testing"
)

func TestRedactorTextIsCaseInsensitive(t *testing.T) {
	r := service.NewRedactor([]string{"token"})
	if got := r.Text("TOKEN=secret"); got == "TOKEN=secret" {
		t.Fatal("sensitive token leaked")
	}
}

func TestRedactorTextRedactsAllCasings(t *testing.T) {
	r := service.NewRedactor([]string{"token"})
	for _, in := range []string{"token=secret", "TOKEN=secret", "Token=secret", "tokEN=secret", "bearer ToKeN xyz"} {
		if got := r.Text(in); got != "[redacted]" {
			t.Fatalf("Text(%q) = %q, want [redacted]", in, got)
		}
	}
}

func TestRedactorTextPreservesOrdinaryText(t *testing.T) {
	r := service.NewRedactor([]string{"token"})
	for _, in := range []string{"hello world", "user logged in", "", "TOKENIZER should not match toKEN"} {
		// "TOKENIZER" must still be redacted because "token" is a substring of it
		// (same behavior as the original case-sensitive Contains). Inputs that do
		// not contain the sensitive term in any casing are returned verbatim.
		if got := r.Text(in); got == "secret" {
			t.Fatalf("Text(%q) unexpectedly redacted non-sensitive input", in)
		}
	}
	// Genuinely clean text passes through unchanged.
	if got := r.Text("ordinary audit message"); got != "ordinary audit message" {
		t.Fatalf("Text(ordinary) = %q, want unchanged", got)
	}
}

func TestRedactorTextWithAuditFlagsChange(t *testing.T) {
	r := service.NewRedactor([]string{"token"})
	if out, changed := r.TextWithAudit("TOKEN=secret"); !changed || out == "TOKEN=secret" {
		t.Fatalf("TextWithAudit(upper) = (%q,%v), want redacted+changed", out, changed)
	}
	if _, changed := r.TextWithAudit("nothing here"); changed {
		t.Fatal("TextWithAudit(clean) should not report a change")
	}
}

func TestPolicyApplyRedactsCaseInsensitive(t *testing.T) {
	p := service.BuildPolicy([]string{"token"})
	for _, in := range []string{"token=secret", "TOKEN=secret", "Token=secret", "tokEN=secret"} {
		if got := p.Apply(in); got == in {
			t.Fatalf("Apply(%q) left sensitive content unredacted: %q", in, got)
		}
	}
	// Non-sensitive content is preserved verbatim.
	if got := p.Apply("ordinary audit message"); got != "ordinary audit message" {
		t.Fatalf("Apply(ordinary) = %q, want unchanged", got)
	}
}

func TestPolicyContainsSensitiveCaseInsensitive(t *testing.T) {
	p := service.BuildPolicy([]string{"token"})
	if !p.ContainsSensitive("Authorization: ToKeN 123") {
		t.Fatal("ContainsSensitive missed mixed-case token")
	}
	if p.ContainsSensitive("no secrets here") {
		t.Fatal("ContainsSensitive false-positive on clean text")
	}
}

func TestPolicyAppliesExactFieldNotSubstring(t *testing.T) {
	// Field-name lookup must stay exact: "ssn" must not match "session".
	p := service.BuildPolicy([]string{"ssn"})
	if p.Applies("session") {
		t.Fatal("Applies matched by substring; field lookup must be exact")
	}
	if !p.Applies("ssn") {
		t.Fatal("Applies did not match the exact field name")
	}
	if !p.Applies("SSN") {
		t.Fatal("Applies should be case-insensitive on field names by default")
	}
}

func TestPolicyCaseSensitiveOptIn(t *testing.T) {
	// When CaseSensitive is explicitly enabled, casing must matter again.
	p := service.RedactionPolicy{Fields: []string{"token"}, Replacement: "[redacted]", CaseSensitive: true}.Normalize()
	if got := p.Apply("TOKEN=secret"); got != "TOKEN=secret" {
		t.Fatalf("case-sensitive policy redacted an upper-case variant: %q", got)
	}
	if got := p.Apply("token=secret"); got == "token=secret" {
		t.Fatal("case-sensitive policy failed to redact exact-case match")
	}
}
