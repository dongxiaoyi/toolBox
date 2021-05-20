# linux操作工具集

# 一、编译：
## 1.指定GOROOT环境变量
> 项目目录下的build目录的绝对路径，如：GOROOT=/data/github.com/dongxiaoyi/toolBox/build

## 2.指定GOPATH环境变量
> 项目目录下的build/gowork目录的绝对路径，如：GOPATH=/data/github.com/dongxiaoyi/toolBox/build/gowork

## 3.其他相关环境变量
- GO111MODULE=on 用于使用go modules做项目依赖管理
- CGO_ENABLED=0 不允许在go代码中调用C代码
- GOOS=linux 编译linux下可用的架构

## 4.其他说明
- go语言的指令路径：项目目录下的build/bin
- 项目依赖包缓存的路径：项目目录下的build/gowork/pkg/mod

## 5.编译示例
```shell
$ cp toolBox.go build/gowork/src/toolBox/
$ cd build/gowork/src/toolBox/
$ GOROOT=/data/github.com/dongxiaoyi/toolBox/build GOPATH=/data/github.com/dongxiaoyi/toolBox/build/gowork GO111MODULE=on CGO_ENABLED=0 GOOS=linux  /data/github.com/dongxiaoyi/toolBox/build/bin/go build
```
说明：
- 编译完成后会在build/gowork/src/toolBox/目录下生成toolBox的二进制文件。

# 二、操作
> 两种操作方式：配置文件 or 命令行参数

```shell
$ ./toolBox --help     
Usage:
  toolBox [command]

Available Commands:
  help        Help about any command
  mod         module

Flags:
  -h, --help   help for toolBox

Use "toolBox [command] --help" for more information about a command.
```

目前支持的模块：

|模块名称|操作|说明|
|---|---|---|
|ntp|server|启动一个本地的ntp server|
|ntp|client|与上述ntp server对应，返回主机与ntp server的时间偏移|
|ntp|netclient|返回与网络ntp server的时间差|
|shell|execute|执行shell指令(一次性显示全部返回值)|
|shell|stream|执行shell指令(流式显示)|
|shell|scp|传递本地文件到远程主机(单文件)|
|telnet|test|测试远程主机的端口是否开放|

## 2.1 配置文件方式操作
> 配置configs/actions.ini文件，如下给出了目前支持的模块的参考配置

`configs/actions.ini`说明：
```ini
[ntp.server]
address = """
172.16.4.129:8888
"""

[ntp.client]
ntp_server = """
172.16.4.129:8888
"""

[ntp.netclient]
ntp_server = """
ntp.aliyun.com
"""

[shell.execute]
ip_port_user_pass_cmd = """
172.16.4.111:22:yanfa:redhat@2020 "echo cmdString"
"""

[shell.stream]
ip_port_user_pass_cmd = """
172.16.4.111:22:yanfa:redhat@2020 "echo cmdString"
"""

[shell.scp]
ip_port_user_pass_src_dest = """
172.16.4.110:22:yanfa:redhat@2020 /tmp/xxx /home/yanfa/xxx
"""

[telnet.test]
ip_port_list = """
172.16.4.129:22 172.16.4.129:16380 172.16.4.129:16381 172.16.4.129:16382
172.16.4.136:22 172.16.4.136:17380
172.16.4.137:22
"""
```

执行操作示例：

```shell
$ ./toolBox mod ntp server
$ ./toolBox mod ntp client
$ ./toolBox mod ntp netclient
$ ./toolBox mod shell execute
$ ./toolBox mod shell stream
$ ./toolBox mod shell scp
$ ./toolBox mod telnet test
```

## 2. 命令行方式直接执行
> 上述方式需要当然，也可以不做配置文件配置，直接命令行执行，示例如：

```shell
$ ./toolBox mod ntp server from=str 172.16.4.128:8888
$ ./toolBox mod ntp client from=str 172.16.4.128:8888
$ ./toolBox mod ntp netclient from=str ntp.aliyun.com
$ ./toolBox mod shell execute from=str 172.16.4.111:22:yanfa:redhat@2020 "echo cmdString"
$ ./toolBox mod shell stream from=str 172.16.4.111:22:yanfa:redhat@2020 "echo cmdString"
$ ./toolBox mod shell scp from=str 172.16.4.110:22:yanfa:redhat@2020 /tmp/xxx /home/yanfa/xxx
$ ./toolBox mod telnet test from=str 172.16.4.129:22 172.16.4.129:16380 172.16.4.129:16381 172.16.4.129:16382
```