# Bug reproduction

The seeded defect is: negative progress chunk changes batch state.

Affected symbol: `internal/domain/masking/model.go::Advance`.

A focused regression test demonstrates the failure before the fix and passes after the fix.
