# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 `golang:1.26.5` 镜像，在当前机器的 `linux/amd64` 容器中验证。标准命令为 `go version`、`go build ./...` 和 `go test ./...`。

## 环境构建与编译

Bug 环境镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。Gold 环境镜像构建成功，容器内 `go build ./...` 和 `go test ./...` 均成功。

## 故障触发步骤

在 Bug 环境执行：

```text
go test -count=1 -run ^TestRejectsTrailingJSON$ ./internal/lab
```

## 实际错误输出

```text
--- FAIL: TestRejectsTrailingJSON (0.00s)
    repro_test.go:18: status=201 body={"id":"obs-000001","site":"wetland","observed_at":"2026-08-18T10:00:00Z","description":"Research team recorded repeated calls beside the marsh reeds.","state":"captured","labels":null,"created_at":"2026-08-19T13:33:57.034840383Z","updated_at":"2026-08-19T13:33:57.034840383Z"}
FAIL
FAIL	example.com/fieldnote-signal-lab/internal/lab	0.011s
FAIL
```

## 期望行为

一次提交只能包含一个完整 JSON 对象；正文存在第二个对象时应返回请求错误，不能创建观测后忽略尾随内容。
