package libs

import (
	"encoding/json"
	"fmt"
	"github.com/dongxiaoyi/toolBox/configs"
	"github.com/dongxiaoyi/toolBox/internal"
	"os"
	"strings"
)

func (t T) NtpServer(m map[string]string) {
	logger := configs.NewLogger(true, true, true, true, false, "console")
	if _, ok := m["cmdContent"]; ok {
		m["address"] = m["cmdContent"]
		delete(m, "cmdContent")
	}
	if m["address"] == "" {
		logger.Error("请填入需要启动的ntp服务地址")
		os.Exit(3)
	}

	internal.NtpServer(strings.TrimSpace(m["address"]))
}

func (t T) NtpClient(m map[string]string) {
	if _, ok := m["cmdContent"]; ok {
		m["ntp_server"] = m["cmdContent"]
		delete(m, "cmdContent")
	}
	if m["ntp_server"] == "" {
		fmt.Println("请填入需要检查的ntp server地址")
		os.Exit(1)
	}

	result := internal.NtpClientRemote(m["ntp_server"])
	js, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println(string(js))
}

func (t T) NtpNetclient(m map[string]string) {
	if _, ok := m["cmdContent"]; ok {
		m["ntp_server"] = m["cmdContent"]
		delete(m, "cmdContent")
	}
	if m["ntp_server"] == "" {
		fmt.Println("请填入需要检查的ntp server地址")
		os.Exit(1)
	}

	result := internal.NtpNetclient(strings.TrimSpace(m["ntp_server"]))

	e, _ := json.Marshal(result)
	fmt.Println(string(e))
}