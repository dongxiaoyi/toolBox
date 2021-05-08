package internal

import (
	netntp "github.com/beevik/ntp"
	ntp "github.com/lixiangyun/go_ntp"
	"time"
)

type NtpMsg struct {
	Code int `json:"code"`
	Message []interface{} `json:"message"`
}

type NtpMsgRemote struct {
	Code int `json:"code"`
	Message interface{} `json:"message"`
}

type NtpNetResult struct {
	Offset int64 `json:"offset"`
	NetDelay int64 `json:"net_delay"`
	Unit string `json:"unit"`
}

// 独立的ntp client服务
func NtpClientRemote(ntpServerStr string) *NtpMsg {
	var msg NtpMsg
	msg.Code = 0

	ntpc := ntp.NewNTPC(ntpServerStr, 1*time.Second)

	// run Delay
	time.Sleep(time.Second)

	/* Initiating synchronization time, waiting for 10 results. */
	// Warn: 源码基础上增加了err返回
	resultAry, err := ntpc.SyncBatch(10)
	if err != nil {
		msg.Code = 2
		msg.Message = []interface{}{float64(0), false, err.Error()}
	} else {
		result := ntp.ResultAverage(resultAry)

		msg.Message = []interface{}{float64(result.Offset.NanoSecond/int64(time.Microsecond)), true, NtpNetResult{
			/* 从本地机发送同步要求到ntp服务器的round trip time */
			Offset: result.Offset.NanoSecond/int64(time.Microsecond),
			/* 主机通过NTP时钟同步与所同步时间源的时间偏移量，单位为毫秒（ms）。offset越接近于0,主机和ntp服务器的时间越接近 */
			NetDelay: result.NetDelay.NanoSecond/int64(time.Microsecond),
			Unit: "us",
		},
		}
	}

	return &msg
}

// 获取外部ntp server的client
func NtpNetclient(ntpServerAddr string) *NtpMsg {
	var msg NtpMsg
	msg.Code = 0

	response, err := netntp.Query(ntpServerAddr)
	if err != nil {
		msg.Code = 2
		msg.Message = []interface{}{float64(0), false, err.Error()}
	} else {
		msg.Message = []interface{}{float64(response.ClockOffset.Microseconds()), true, NtpNetResult{
			/* 从本地机发送同步要求到ntp服务器的round trip time */
			Offset: response.ClockOffset.Microseconds(),
			/* 主机通过NTP时钟同步与所同步时间源的时间偏移量，单位为毫秒（ms）。offset越接近于0,主机和ntp服务器的时间越接近 */
			NetDelay: response.RootDelay.Microseconds(),
			Unit: "us",
		},
		}
	}

	return &msg
}

// 独立的ntp server服务
func NtpServer(addr string) {
	ntps := ntp.NewNTPS(addr)
	// Start the service, then coroutine process is created in the background.
	err := ntps.Start()
	if err != nil {
		panic(err)
	}
	for {
	}
}