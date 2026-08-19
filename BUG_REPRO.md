# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 `golang:1.26.5` 镜像，在当前机器的 `linux/amd64` 容器中验证。标准命令为 `go version`、`go build ./...` 和 `go test ./...`。

## 环境构建与编译

Bug 环境镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。Gold 环境镜像构建成功，容器内 `go build ./...` 和 `go test ./...` 均成功。

## 故障触发步骤

在 Bug 环境执行：

```text
go test -count=1 -run ^TestListTotalCountsAllMatches$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestListTotalCountsAllMatches (0.00s)
    repro_test.go:35: page={Items:[{ID:obs-000001 Site:wetland ObservedAt:2026-08-19 12:32:21 +0000 UTC Description:Research team recorded repeated calls beside the marsh reeds. State:captured Labels:[] Review:<nil> CreatedAt:2026-08-19 13:32:21.56219741 +0000 UTC UpdatedAt:2026-08-19 13:32:21.56219741 +0000 UTC}] Total:1}
FAIL
FAIL	example.com/fieldnote-signal-lab/internal/lab	0.011s
FAIL
```

## 期望行为

页面只展示一条当前页记录时，总记录数仍应为全部匹配的两条记录。
