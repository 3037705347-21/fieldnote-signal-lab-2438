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
go test -count=1 -run ^TestObservationReadsAreIsolated$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestObservationReadsAreIsolated (0.00s)
    repro_test.go:33: stored labels were mutated through a read: [pollinator nesting]
FAIL
FAIL    example.com/fieldnote-signal-lab/internal/lab    0.250s
FAIL
```

## 期望行为

调用方修改读取结果后，已保存的观测、标签和复核信息都不应被改变，后续查询应返回原始数据。
