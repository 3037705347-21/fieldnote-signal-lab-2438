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

使用 `golang:1.26.5` 构建 Docker 镜像。镜像构建成功，容器内 `go build ./...` 成功；并发复现使用 Linux 容器的 `-race` 检查。

## 故障触发步骤

在 bug 环境中执行：

```text
go test -race -count=1 -run ^TestConcurrentLabelsAreMerged$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestConcurrentLabelsAreMerged (0.00s)
    repro_test.go:78: labels=[nesting]
FAIL
FAIL    example.com/fieldnote-signal-lab/internal/lab    0.016s
FAIL
```

## 期望行为

两个研究员同时为同一条观测提交不同标签后，两种标签都应保留，并继续满足标签去重、排序和状态流转要求。
