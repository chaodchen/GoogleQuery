## fish666 project
基于fyne框架实现的跨平台谷歌收录查询工具

### 环境
- golang:1.20
- fyne:2.3.5

### 编译
First
```shell
go install fyne.io/fyne/v2/cmd/fyne@latest
```
ForWindows:
```shell
fyne package -os windows -icon icon.png
```
ForMac:
```shell
fyne package -os darwin -icon icon.png
```
ForLinux:
```shell
fyne package -os linux -icon icon.png
```
### 运行
```go
go mod tidy
go run main.go
```
