package internal

import (
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// 执行远程shell指令
func ShellExecuteRemote(ip, port, user string, password interface{}, cmdStr string) (string, error) {
	var config *ssh.ClientConfig
	if password != nil {
		// 账号密码登录
		config = &ssh.ClientConfig{
			Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            user,
			Auth: []ssh.AuthMethod{ssh.Password(password.(string))},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
	} else {
		// 免密登录
		userDir := fmt.Sprintf(`/home/%s`, user)
		if user == "root" {
			userDir = "/root"
		}

		pemKeyPath := path.Join(userDir, ".ssh/id_rsa")
		pemBytes, err := ioutil.ReadFile(pemKeyPath)
		if err != nil {
			return ip +"获取私钥失败", err
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return ip +"解析私钥失败", err
		}
		config = &ssh.ClientConfig{
			Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            user,
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
	}

	sshClient, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		return ip +"ssh连接失败", err
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		return ip +"ssh session连接失败", err
	}
	defer session.Close()
	//执行远程命令
	combo,err := session.CombinedOutput(cmdStr)
	if err != nil {
		return ip +"远程指令["+cmdStr+"]操作失败", err
	}
	return string(combo), nil
}

// 执行远程shell指令(流式输出结果)
func ShellExecuteRemoteStream(ip, port, user string, password interface{}, cmdStr string) error {
	var config *ssh.ClientConfig
	if password != nil {
		config = &ssh.ClientConfig{
			Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            user,
			Auth: []ssh.AuthMethod{ssh.Password(password.(string))},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
	} else {
		// 免密登录
		userDir := fmt.Sprintf(`/home/%s`, user)
		if user == "root" {
			userDir = "/root"
		}

		pemKeyPath := path.Join(userDir, ".ssh/id_rsa")
		pemBytes, err := ioutil.ReadFile(pemKeyPath)
		if err != nil {
			return err
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return err
		}
		config = &ssh.ClientConfig{
			Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            user,
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
	}

	sshClient, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	session.Stderr = session.Stdout
	//执行远程命令
	//_, err = session.CombinedOutput(cmdStr)
	err = session.Start(cmdStr)
	if err != nil {
		return err
	}
	for {
		tmp := make([]byte, 1024)
		_, err = stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
	err = session.Wait()
	if err != nil {
		return err
	}
	return nil
}

func ShellScpRemote(ip, port, user string, password interface{}, src, dest string) (string, error) {
	var config *ssh.ClientConfig
	if password != nil {
		config = &ssh.ClientConfig{
			Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            user,
			Auth: []ssh.AuthMethod{ssh.Password(password.(string))},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
	} else {
		// 免密登录
		userDir := fmt.Sprintf(`/home/%s`, user)
		if user == "root" {
			userDir = "/root"
		}

		pemKeyPath := path.Join(userDir, ".ssh/id_rsa")
		pemBytes, err := ioutil.ReadFile(pemKeyPath)
		if err != nil {
			return ip +"获取私钥失败", err
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return ip +"解析私钥失败", err
		}
		config = &ssh.ClientConfig{
			Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            user,
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
	}

	sshClient, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		msgSshClient := fmt.Sprintf(`%s 创建ssh client 失败`, ip)
		return msgSshClient, err
	} else {
		defer sshClient.Close()
	}

	client, err := scp.NewClientBySSH(sshClient)
	if err != nil {
		msgScpClient := fmt.Sprintf(`scp client 失败, %v`, err)
		return msgScpClient, err
	} else {
		defer client.Close()
	}

	f, _ := os.Open(src)

	err = client.CopyFile(f, dest, "0755")

	if err != nil {
		msgCopyDest := fmt.Sprintf(`Error while copying file , Maybe the path does not exist, %v`, err)
		return msgCopyDest, err
	} else {
		f.Close()
	}

	msgMessageSuccess := fmt.Sprintf(`ip [%s] scp文件%v成功！`, ip , dest)
	return msgMessageSuccess, nil
}