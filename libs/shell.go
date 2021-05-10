package libs

import (
	"errors"
	"fmt"
	"github.com/dongxiaoyi/toolBox/pkg"
	"os"
	"strings"
)

/*
Example：
$ ./toolBox mod shell execute from=str 172.16.4.111:22:yanfa:redhat@2020 "if [ ! -f "/data/tmp/rbac.yaml" ]; then echo xxxx; fi"
or
$ ./toolBox mod shell execute from=str 172.16.4.111:22:yanfa "if [ ! -f "/data/tmp/rbac.yaml" ]; then echo xxxx; fi"

说明：
- 不输入密码的时候 -> ip:port:user;输入密码的时候 -> ip:port:user:password
*/
func (t T) ShellExecute(m map[string]string) {
	if _, ok := m["cmdContent"]; ok {
		m["ip_port_user_pass_cmd"] = m["cmdContent"]
		delete(m, "cmdContent")
	}

	if m["ip_port_user_pass_cmd"] == "" {
		fmt.Println("请填入需要远程客户端地址、需要执行的shell指令")
		os.Exit(1)
	}

	ipPortUserPassCmdListStr := strings.TrimSpace(m["ip_port_user_pass_cmd"])
	ipPortUserPassCmdSlice := strings.Fields(ipPortUserPassCmdListStr)

	if len(ipPortUserPassCmdSlice) < 2 {
		fmt.Println("远程客户端地址、需要执行的shell指令 必须全部填写")
		os.Exit(2)
	}

	ipPortUserPassSlice := strings.Split(ipPortUserPassCmdSlice[0], ":")
	cmdClice := ipPortUserPassCmdSlice[1:]
	cmdStr := strings.Join(cmdClice, " ")

	var result string
	var err error

	if len(ipPortUserPassSlice) == 4 {
		result, err = pkg.ShellExecuteRemote(ipPortUserPassSlice[0], ipPortUserPassSlice[1], ipPortUserPassSlice[2], ipPortUserPassSlice[3], cmdStr)
	} else if len(ipPortUserPassSlice) == 3 {
		result, err = pkg.ShellExecuteRemote(ipPortUserPassSlice[0], ipPortUserPassSlice[1], ipPortUserPassSlice[2], nil, cmdStr)
	} else {
		fmt.Println(errors.New("请输入远程地址的ip、port、user、[password]"))
		os.Exit(3)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	fmt.Println(result)
}

/*
Example：
$ ./toolBox mod shell stream from=str 172.16.4.111:22:yanfa:redhat@2020 "if [ ! -f "/data/tmp/rbac.yaml" ]; then echo xxxx; fi"
or
$ ./toolBox mod shell stream from=str 172.16.4.111:22:yanfa "if [ ! -f "/data/tmp/rbac.yaml" ]; then echo xxxx; fi"

说明：
- 不输入密码的时候 -> ip:port:user;输入密码的时候 -> ip:port:user:password
*/
func (t T) ShellStream(m map[string]string) {
	if _, ok := m["cmdContent"]; ok {
		m["ip_port_user_pass_cmd"] = m["cmdContent"]
		delete(m, "cmdContent")
	}

	if m["ip_port_user_pass_cmd"] == "" {
		fmt.Println("请填入需要远程客户端地址、需要执行的shell指令")
		os.Exit(1)
	}

	ipPortUserPassCmdListStr := strings.TrimSpace(m["ip_port_user_pass_cmd"])
	ipPortUserPassCmdSlice := strings.Fields(ipPortUserPassCmdListStr)

	if len(ipPortUserPassCmdSlice) < 2 {
		fmt.Println("远程客户端地址、需要执行的shell指令 必须全部填写")
		os.Exit(2)
	}

	ipPortUserPassSlice := strings.Split(ipPortUserPassCmdSlice[0], ":")
	cmdClice := ipPortUserPassCmdSlice[1:]
	cmdStr := strings.Join(cmdClice, " ")
	fmt.Println(cmdStr)

	var err error

	if len(ipPortUserPassSlice) == 4 {
		err = pkg.ShellExecuteRemoteStream(ipPortUserPassSlice[0], ipPortUserPassSlice[1], ipPortUserPassSlice[2], ipPortUserPassSlice[3], cmdStr)
	} else if len(ipPortUserPassSlice) == 3 {
		err = pkg.ShellExecuteRemoteStream(ipPortUserPassSlice[0], ipPortUserPassSlice[1], ipPortUserPassSlice[2], nil, cmdStr)
	} else {
		fmt.Println(errors.New("请输入远程地址的ip、port、user、[password]"))
		os.Exit(3)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
}


/*
Example:
$ ./toolBox mod scp 172.16.4.110:22:yanfa:redhat@2020 /tmp/xxx /home/yanfa/xxx
or
$ ./toolBox mod scp 172.16.4.110:22:yanfa /tmp/xxx /home/yanfa/xxx

说明：
- 不输入密码的时候 -> ip:port:user;输入密码的时候 -> ip:port:user:password
*/
func (t T) ShellScp(m map[string]string) {
	if _, ok := m["cmdContent"]; ok {
		m["ip_port_user_pass_src_dest"] = m["cmdContent"]
		delete(m, "cmdContent")
	}
	if m["ip_port_user_pass_src_dest"] == "" {
		msgMessage := fmt.Sprintf(`请填入需要远程客户端地址、需要推送的文件源地址、目标地址`)
		fmt.Println(msgMessage)
		os.Exit(1)
	}

	ipPortUserPassSrcDestStr := strings.TrimSpace(m["ip_port_user_pass_src_dest"])
	ipPortUserPassSrcDestSlice := strings.Fields(ipPortUserPassSrcDestStr)

	if len(ipPortUserPassSrcDestSlice) < 3 {
		msgMessage := fmt.Sprintf(`远程客户端地址、需要推送的文件源地址、目标地址 必须全部填写`)
		fmt.Println(msgMessage)
		os.Exit(2)
	}
	ipPortUserPassSlice := strings.Split(ipPortUserPassSrcDestSlice[0], ":")

	var msg string
	var err error
	if len(ipPortUserPassSlice) == 4 {
		msg, err = pkg.ShellScpRemote(ipPortUserPassSlice[0], ipPortUserPassSlice[1], ipPortUserPassSlice[2],
			ipPortUserPassSlice[3], ipPortUserPassSrcDestSlice[1], ipPortUserPassSrcDestSlice[2])
	} else if len(ipPortUserPassSlice) == 3 {
		msg, err = pkg.ShellScpRemote(ipPortUserPassSlice[0], ipPortUserPassSlice[1], ipPortUserPassSlice[2],
			nil, ipPortUserPassSrcDestSlice[1], ipPortUserPassSrcDestSlice[2])
	} else {
		fmt.Println(errors.New("请输入远程地址的ip、port、user、[password]"))
		os.Exit(3)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	} else {
		fmt.Println(msg)
	}
}
