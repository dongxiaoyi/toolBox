package libs

import (
	"fmt"
	"github.com/dongxiaoyi/toolBox/pkg"
	"os"
	"strings"
	"time"
)

/*
Example:
$ ./op-tools mod telnet test from=str 172.16.4.128:22 172.16.4.129:22
 */
func (t T) TelnetTest(m map[string]string) {
	// ip:port清单为使用空格间隔的字符串
	if _, ok := m["cmdContent"]; ok {
		m["ip_port_list"] = m["cmdContent"]
		delete(m, "cmdContent")
	}
	if m["ip_port_list"] == "" {
		fmt.Println("请填入待检测ip:port")
		os.Exit(2)
	}

	ipPortListStr := strings.TrimSpace(m["ip_port_list"])
	ipPortSlice := strings.Fields(ipPortListStr)

	result := make(chan []byte)
	go pkg.TelnetConnTest(ipPortSlice, result)

	select {
	case r := <- result:
		fmt.Println(string(r))
		break
	case <- time.After(10 * time.Second):
		fmt.Println("操作超时")
		break
	}
}
