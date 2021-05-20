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


## 1. 模块配置
> 配置configs/actions.ini文件

`configs/actions.ini`说明：
```ini
[mod_name.mod_action]
mod_args = arg1
mod_multiline_args = """
mod_multiline_arg1
mod_multiline_arg2
mod_multiline_arg3
"""
```

## 2. 当前的支持的mod及action

### 2.1 telnet模块
#### 2.2 telnet模块action：test
configs/actions.ini配置方法参考：
```ini
[telnet.test]
ip_port_list = """
172.16.4.129:22 172.16.4.129:16380 172.16.4.129:16381 172.16.4.129:16382
172.16.4.136:22 172.16.4.136:17380
172.16.4.137:22
"""
```
执行操作：

```shell
$ ./op-tools mod telnet test
```

当然，也可以不做配置文件配置，直接命令行执行，如：

```shell
$ ./op-tools mod telnet test from=str 172.16.4.128:22 172.16.4.129:22
```

其他：
- closestmatch要使用github最新的代码，不要使用go mod的版本，go mod版本的代码有bug

BUG:超时后不能统计到结果。