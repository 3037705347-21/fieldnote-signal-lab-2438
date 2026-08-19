# 修复前故障复现（Docker）

## 项目与标准命令

- 平台：Docker `linux/amd64`
- 镜像构建：`docker build --progress=plain -f benzhi.Dockerfile -t go-fieldnote-signal-lab__019-bug:20260819 .`
- 容器验证：`docker run --rm go-fieldnote-signal-lab__019-bug:20260819 sh -c "go test ./..."`

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
--- FAIL: TestReviewCanonicalizesVerdictForScoring (0.00s)
    review_verdict_test.go:32: review=&{Verdict:Confirmed Confidence:0.9 Notes:Independent audio confirms the field note. ReviewedAt:2026-08-19 14:26:23.072535003 +0000 UTC}
FAIL
FAIL	example.com/fieldnote-signal-lab/internal/lab	0.011s
FAIL
```

## 期望行为

带空格或大小写变体的正常复核结论应保存为统一值，并按确认结论参与摘要评分。
