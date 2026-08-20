# Bug reproduction

The seeded defect is: restricted fields without scope are silently accepted.

Affected symbol: `internal/application/service.go::Preview`.

A focused regression test demonstrates the failure before the fix and passes after the fix.
