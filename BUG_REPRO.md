# 修复前故障复现（Docker）

## 项目与标准命令

平台：`linux/amd64`

标准命令：

```text
go version
go build ./...
go test ./...
```

## 环境构建与编译

使用 `golang:1.26.5` 构建 Docker 镜像。镜像构建成功，容器内 `go build ./...` 成功。

## 故障触发步骤

在 bug 环境中执行：

```text
go test -count=1 -run ^TestListTotalIgnoresLimit$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestListTotalIgnoresLimit (0.00s)
    repro_test.go:35: page={Items:[{ID:obs-000001 Site:wetland ObservedAt:2026-08-19 09:08:55 +0000 UTC Description:Research team recorded repeated calls beside the marsh reeds. State:captured Labels:[] Review:<nil> CreatedAt:2026-08-19 10:08:55.0099135 +0000 UTC UpdatedAt:2026-08-19 10:08:55.0099135 +0000 UTC} {ID:obs-000002 Site:wetland ObservedAt:2026-08-19 08:08:55 +0000 UTC Description:Research team recorded repeated calls beside the marsh reeds. State:captured Labels:[] Review:<nil> CreatedAt:2026-08-19 10:08:55.0099135 +0000 UTC UpdatedAt:2026-08-19 10:08:55.0099135 +0000 UTC}] Total:2}
FAIL
FAIL    example.com/fieldnote-signal-lab/internal/lab    0.165s
FAIL
```

## 期望行为

列表只返回当前页记录，但匹配总数应统计所有符合筛选条件的记录，而不是当前页的记录数。
