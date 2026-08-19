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
go test -count=1 -run ^TestListTagFilterTrimsWhitespace$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestListTagFilterTrimsWhitespace (0.00s)
    repro_test.go:26: page={Items:[] Total:0}
FAIL
FAIL    example.com/fieldnote-signal-lab/internal/lab    0.012s
FAIL
```

## 期望行为

按受支持的标签筛选观测记录时，即使输入值带有首尾空白，也应返回匹配记录，并保持标签大小写规范化、去重排序及其他筛选条件的现有行为。
