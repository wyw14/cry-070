# Bug reproduction

The seeded defect is: batch runner ignores a cancellation error at persistence boundary.

Affected symbol: `internal/application/service.go::RunBatch`.

A focused regression test demonstrates the failure before the fix and passes after the fix.
