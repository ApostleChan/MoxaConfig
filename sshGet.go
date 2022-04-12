package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func NewSftpclient(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
		//ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey)
	}
	// connet to ssh
	fmt.Println(user, password)
	addr = fmt.Sprintf("%s:%d", host, port)
	fmt.Println(addr + "connect success")
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func SshInteractive(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	answers = make([]string, len(questions))
	// The second parameter is unused
	for n, _ := range questions {
		answers[n] = "your_password"
	}

	return answers, nil
}

func DownnloadFile(sftpClient *sftp.Client, localPath string, remoteFilePath string) (bool, error) {
	p, _ := os.Stat(localPath)
	if p == nil {
		fmt.Println("localpath is not exist")
		err := os.Mkdir(localPath, 0777)
		if err != nil {
			return false, err
		}
		fmt.Println("create loaclpath success")
	}
	var remoteFileName = path.Base(remoteFilePath)
	//fmt.Println(remoteFileName)
	dstFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		fmt.Println("sftpClient.Open error : ", remoteFilePath)
		return false, err
	}
	defer dstFile.Close()
	localFileName, err := os.Create(path.Join(localPath, remoteFileName))
	if err != nil {
		fmt.Println("os.Create failed", localFileName)
		return false, err
	}
	_, err = dstFile.WriteTo(localFileName)
	if err != nil {
		fmt.Println("download success")
		return true, nil
	}
	return false, err

}

func UploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) (bool, error) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		return false, err
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		return false, err
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		return false, err
	}
	_, err = dstFile.Write(ff)
	if err != nil {
		return false, err
	}
	fmt.Println(localFilePath + " copy file to remote server finished!")
	return true, nil
}

func UploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("read dir list fail ", err)
	}
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			sftpClient.Mkdir(remoteFilePath)
			UploadDirectory(sftpClient, localFilePath, remoteFilePath)
		} else {
			UploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
		}
	}
	fmt.Println(localPath + " copy directory to remote server finished!")
}

func CreateNewDir(dirName string) {
	current, _ := os.Getwd()
	os.Mkdir(path.Join(current, dirName), 0777)
}
