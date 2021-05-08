package internal

import (
	"encoding/json"
	tc "github.com/reiver/go-telnet"
	"strings"
	"sync"
	"time"
)

const timeout = 8 * time.Second

type TelnetTestMsg struct {
	Result bool `json:"result"`
	Message string `json:"message"`
}

var (
	 // FIXME: 不应该使用全局变量
	lock sync.Mutex
)

/*
Example:
$ ./op-tools mod telnet test from=str 172.16.4.128:22 172.16.4.129:22
*/
func TelnetConnTest(ipPortSlice []string, result chan []byte) {
	// ip:port清单为使用空格间隔的字符串
	counterResult := make(map[string]TelnetTestMsg)
	tcPoolChan := make(chan TelnetTestMsg, len(ipPortSlice))
	for i:=0; i < len(ipPortSlice); i++ {
		ipPortList := strings.Split(ipPortSlice[i], ":")
		if len(ipPortList) == 1 {
			ipPortList = append(ipPortList, "None")
		}
		resultMsgInit := TelnetTestMsg{
			Result: false,
			Message: "Maybe connetion time out.",
		}
		counterResult[ipPortSlice[i]] = resultMsgInit
	}

	var wg sync.WaitGroup
	for i:=0; i < len(ipPortSlice); i++ {
		//tc.StandardCaller
		wg.Add(1)
		go telnetClientConnectionTest(tcPoolChan, ipPortSlice[i], &counterResult)
	}

	counterSucc := 0
	counterFail := 0

	loop:
	for {
		select {
		case t := <- tcPoolChan:
			if t.Result {
				counterSucc += 1
			} else {
				counterFail += 1
			}

			if counterFail + counterSucc == len(ipPortSlice) {
				resultByte, _ := json.Marshal(counterResult)
				result <- resultByte
				break loop
			}
		case <-time.After(timeout):
			resultByte, _ := json.Marshal(counterResult)
			result <- resultByte
			break loop
		default:
			continue
		}
	}
	wg.Wait()
}

func telnetClientConnectionTest(msg chan TelnetTestMsg, ipPort string, counterResult *map[string]TelnetTestMsg) {
	dail, err := tc.DialTo(ipPort)

	ipPortSlice := strings.Split(ipPort, ":")
	if len(ipPortSlice) == 1 {
		ipPortSlice = append(ipPortSlice, "None")
	}

	lock.Lock()
	(*counterResult)[ipPort] = TelnetTestMsg{
		Result: true,
		Message: "Connction success.",
	}

	if err != nil {
		(*counterResult)[ipPort] = TelnetTestMsg{
			Result: false,
			Message: "Connction error.",
		}
	} else {
		dail.Close()
	}
	lock.Unlock()
	msg <- (*counterResult)[ipPort]
}
