# Bug reproduction

The seeded defect is: succeeded batch can be cancelled.

Affected symbol: `internal/domain/masking/model.go::Cancel`.

A focused regression test demonstrates the failure before the fix and passes after the fix.
