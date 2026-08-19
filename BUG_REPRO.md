# BUG_REPRO

## 项目与标准命令

Fieldnote Signal Lab 的观测列表接口应按受控标签查询记录。

## 环境构建与编译

```bash
go build ./...
```

编译通过。

## 故障触发步骤

```bash
go test -count=1 -run ^TestListNormalizesTagQuery$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestListNormalizesTagQuery
    repro_test.go:33: page={Items:[] Total:0}
FAIL
```

## 期望行为

已经附加 acoustic 标签的观测使用 `Acoustic` 查询时仍应返回对应记录。
