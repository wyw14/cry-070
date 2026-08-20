# Bug reproduction

The seeded defect is: two batches with the same idempotency key are both created.

Affected symbol: `internal/repository/memory/store.go::CreateBatch`.

A focused regression test demonstrates the failure before the fix and passes after the fix.
