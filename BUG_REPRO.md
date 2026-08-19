# BUG_REPRO

## 项目与标准命令

Fieldnote Signal Lab 支持在复核后继续补充受控标签。

## 环境构建与编译

```bash
go build ./...
```

编译通过。

## 故障触发步骤

```bash
go test -count=1 -run ^TestAddingLabelKeepsReviewedObservationReviewed$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestAddingLabelKeepsReviewedObservationReviewed
    repro_test.go:27: state=labeled labels=[acoustic habitat nesting weather]
FAIL
```

## 期望行为

已完成复核的观测补充标签后应继续保持 reviewed 状态。
