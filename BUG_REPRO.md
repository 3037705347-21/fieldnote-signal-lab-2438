# BUG_REPRO

## 项目与标准命令

Fieldnote Signal Lab 的内存仓库应把读取结果与已保存状态隔离。

## 环境构建与编译

```bash
go build ./...
```

编译通过。

## 故障触发步骤

```bash
go test -count=1 -run ^TestReadingObservationDoesNotAliasStoredLabels$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestReadingObservationDoesNotAliasStoredLabels
    repro_test.go:20: stored labels were mutated: [weather]
FAIL
```

## 期望行为

调用方修改一次读取结果的标签，不应改变后续读取到的已保存标签。
