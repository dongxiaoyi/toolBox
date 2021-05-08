package libs

import (
	"fmt"
	"github.com/dongxiaoyi/toolBox/internal"
	"os"
	"strings"
)

// ./toolBox mod shell execute from=str 172.16.4.111:22:yanfa:redhat@2020 "if [ ! -f "/data/tmp/rbac.yaml" ]; then echo xxxx; fi"
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

	result, err := internal.ShellExecuteRemote(ipPortUserPassSlice[0], ipPortUserPassSlice[1], ipPortUserPassSlice[2], ipPortUserPassSlice[3], cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(result)
}

// ./toolBox mod shell stream from=str 172.16.4.111:22:yanfa:redhat@2020 "if [ ! -f "/data/tmp/rbac.yaml" ]; then echo xxxx; fi"
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

	err := internal.ShellExecuteRemoteStream(ipPortUserPassSlice[0], ipPortUserPassSlice[1], ipPortUserPassSlice[2], ipPortUserPassSlice[3], cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

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

	msg, err := internal.ShellScpRemote(ipPortUserPassSrcDestSlice[0], ipPortUserPassSrcDestSlice[1], ipPortUserPassSrcDestSlice[2],
		ipPortUserPassSrcDestSlice[3], ipPortUserPassSrcDestSlice[4], ipPortUserPassSrcDestSlice[5])
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	} else {
		fmt.Println(msg)
	}
}
