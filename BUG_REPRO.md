# BUG_REPRO

## 项目与标准命令

Fieldnote Signal Lab 提供按地点筛选并限制返回数量的观测列表接口。

## 环境构建与编译

```bash
go build ./...
```

编译通过。

## 故障触发步骤

```bash
go test -count=1 -run ^TestListAppliesLimitAfterFiltering$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestListAppliesLimitAfterFiltering
    repro_test.go:31: page={Items:[] Total:0}
FAIL
```

## 期望行为

前序记录不匹配地点筛选时，受限查询仍应返回后续匹配的森林观测。
