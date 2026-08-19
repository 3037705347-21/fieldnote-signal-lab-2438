# BUG_REPRO

## 项目与标准命令

Fieldnote Signal Lab 会根据观测标签、描述和复核结论生成信号报告。

## 环境构建与编译

```bash
go build ./...
```

编译通过。

## 故障触发步骤

```bash
go test -count=1 -run ^TestRejectedReviewDoesNotIncreaseReportScore$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestRejectedReviewDoesNotIncreaseReportScore
    repro_test.go:38: score=54 reasons=[description contains sufficient field detail controlled tags were attached review confidence was incorporated]
FAIL
```

## 期望行为

被拒绝的复核不应增加报告的正向分值，示例报告分数应为 49。
