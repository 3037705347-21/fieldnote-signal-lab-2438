# 修复前故障复现（Docker）

## 项目与标准命令

- 平台：Docker `linux/amd64`
- 镜像构建：`docker build --progress=plain -f benzhi.Dockerfile -t go-fieldnote-signal-lab__018-bug:20260819 .`
- 容器验证：`docker run --rm go-fieldnote-signal-lab__018-bug:20260819 sh -c "go test ./..."`

## 环境构建与编译

Bug 环境镜像构建成功，容器内执行 `go version` 和 `go build ./...` 均通过。

## 故障触发步骤

在容器内执行：

```text
go test ./...
```

## 实际错误输出

```text
?   	example.com/fieldnote-signal-lab/cmd/api	[no test files]
--- FAIL: TestRejectsInvalidListLimit (0.00s)
    invalid_limit_test.go:18: status=200, want 400
FAIL
FAIL	example.com/fieldnote-signal-lab/internal/lab	0.012s
FAIL
```

## 期望行为

列表分页参数不是数字时应明确提示输入无效，而不是按不限量成功返回全部记录。
