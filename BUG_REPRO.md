# Bug reproduction

The seeded defect is: log redaction is case-sensitive and leaks uppercase tokens.

Affected symbol: `internal/service/redaction.go::Text`.

A focused regression test demonstrates the failure before the fix and passes after the fix.
