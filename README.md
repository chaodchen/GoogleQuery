## fish666 project
基于fyne框架实现的跨平台谷歌收录查询工具

### 需求
- [ ] 代理功能自定义 目前写死不方便配置
- [ ] 表格过长影响查看、排序功能 字段升序倒序
- [x] 导出功能 导出完成弹个框
- [ ] logs  显示当前查询进度、当前查询域名
- [ ] 可以适当延迟个 2-3秒 不追求速度 追求准确

### 环境
- golang:1.20
- fyne:2.3.5

### 打包
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
### 编译 & 运行
```shell
go mod tidy
go build
./fish666
```


