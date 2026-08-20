# Bug reproduction

The seeded defect is: row diff silently drops target-only rows.

Affected symbol: `internal/domain/masking/diff.go::CompareRows`.

A focused regression test demonstrates the failure before the fix and passes after the fix.
