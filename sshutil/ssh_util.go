package sshutil

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

var waitgroup sync.WaitGroup

//递归下载服务器文件夹
func ListDownloadFile(sftpClient *sftp.Client, dir, savePath string) {
	if !exists(savePath) {
		os.Mkdir(savePath, os.ModePerm)
	}
	fileInfos, err := sftpClient.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, fileInfo := range fileInfos {
			fmt.Println(fileInfo.Name())
			if fileInfo.IsDir() {
				os.Mkdir(savePath+"/"+fileInfo.Name(), os.ModePerm)
				ListDownloadFile(sftpClient, dir+"/"+fileInfo.Name(), savePath+"/"+fileInfo.Name())
			} else {
				waitgroup.Add(1)
				go DownloadFile(sftpClient, dir, savePath, fileInfo.Name())
			}
		}
	}
	waitgroup.Wait()
}

//下载服务器文件
func DownloadFile(sftpClient *sftp.Client, dir, savePath, name string) {
	defer waitgroup.Done()
	srcFile, err := sftpClient.Open(dir + "/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	var savePathFile = path.Join(savePath, name)
	if exists(savePathFile) {
		fmt.Println("file exists!")
		return
	}
	dstFile, err := os.Create(savePathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		log.Fatal(err)
	}

	fmt.Println("copy file from remote server finished!")
}

//判断文件是否存在
func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func SftpConnectPass(user, password, host string, port int) (sftpClient *sftp.Client, err error) { //参数: 远程服务器用户名, 密码, ip, 端口
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := host + ":" + strconv.Itoa(port)
	sshClient, err := ssh.Dial("tcp", addr, clientConfig) //连接ssh
	if err != nil {
		fmt.Println("连接ssh失败", err)
		return
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil { //创建客户端
		fmt.Println("创建客户端失败", err)
		return
	}

	return
}

func SftpConnectSecret(user, keyPath, host string, port int) (sftpClient *sftp.Client, err error) {
	auth := make([]ssh.AuthMethod, 0)

	privateKeyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	auth = append(auth, ssh.PublicKeys(key))

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := host + ":" + strconv.Itoa(port)
	sshClient, err := ssh.Dial("tcp", addr, clientConfig) //连接ssh
	if err != nil {
		fmt.Println("连接ssh失败", err)
		return
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil { //创建客户端
		fmt.Println("创建客户端失败", err)
		return
	}

	return
}

func SftpConnectSecretPass(user, keyPath, keyPass, host string, port int) (sftpClient *sftp.Client, err error) {
	auth := make([]ssh.AuthMethod, 0)

	privateKeyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ssh.ParsePrivateKeyWithPassphrase(privateKeyBytes, []byte(keyPass))
	if err != nil {
		log.Fatal(err)
	}

	auth = append(auth, ssh.PublicKeys(key))

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := host + ":" + strconv.Itoa(port)
	sshClient, err := ssh.Dial("tcp", addr, clientConfig) //连接ssh
	if err != nil {
		fmt.Println("连接ssh失败", err)
		return
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil { //创建客户端
		fmt.Println("创建客户端失败", err)
		return
	}

	return
}
