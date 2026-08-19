# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 `golang:1.26.5` 镜像，在当前机器的 `linux/amd64` 容器中验证。标准命令为 `go version`、`go build ./...` 和 `go test ./...`。

## 环境构建与编译

Bug 环境镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。Gold 环境镜像构建成功，容器内 `go build ./...` 和 `go test ./...` 均成功。

## 故障触发步骤

在 Bug 环境执行：

```text
go test -count=1 -run ^TestReportIncludesFullLabelContribution$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestReportIncludesFullLabelContribution (0.00s)
    repro_test.go:24: report={ObservationID:demo-001 Score:52 Band:developing Reasons:[description contains sufficient field detail controlled tags were attached review confidence was incorporated] GeneratedAt:2026-08-19 13:35:30.020192721 +0000 UTC}
FAIL
FAIL	example.com/fieldnote-signal-lab/internal/lab	0.011s
FAIL
```

## 期望行为

包含三个受支持标签和高置信度复核的 `demo-001` 摘要分数应为 79。
