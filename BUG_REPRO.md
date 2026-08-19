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
go test -count=1 -run ^TestReviewNormalizesVerdict$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestReviewNormalizesVerdict (0.00s)
    repro_test.go:23: status=400 body={"code":"invalid_request","message":"verdict: is unsupported"}
FAIL
FAIL    example.com/fieldnote-signal-lab/internal/lab    0.148s
FAIL
```

## 期望行为

提交受控观测结论时，结论大小写不同或带有首尾空白不应被错误拒绝，并应继续保存为规范化的有效结论。
