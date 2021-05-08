# linux操作工具集

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