# 离线数据脱敏流水线治理平台

平台登记本地样例数据源、字段敏感级别、脱敏策略和批次，必须先预演确认再执行；源表永不原地覆盖。所有通知、快照和报告适配器均为本地实现。

运行：`go mod tidy`、`go test ./...`、`go test -race ./...`、`go vet ./...`、`go run ./cmd/server`。服务默认监听 `:8092`，API 前缀 `/api/v1`，迁移在 `migrations/001_init.sql`，前端在 `web/`。
