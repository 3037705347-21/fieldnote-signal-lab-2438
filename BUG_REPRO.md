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
go test -count=1 -run ^TestListSiteFilterIsCaseInsensitive$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestListSiteFilterIsCaseInsensitive (0.00s)
    repro_test.go:31: page={Items:[] Total:0}
FAIL
FAIL    example.com/fieldnote-signal-lab/internal/lab    0.139s
FAIL
```

## 期望行为

站点名称仅大小写不同但表示同一站点时，筛选结果应包含该站点的观测记录。
