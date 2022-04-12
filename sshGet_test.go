package main

import (
	"os"
	"path"
	"testing"
)

func TestTNewSftpclient(t *testing.T) {
	currentPath, _ := os.Getwd()
	sftpObj, err := NewSftpclient("root", "root", "192.168.1.7", 22)
	if err != nil {
		t.Log(err)
	}
	result, _ := DownnloadFile(sftpObj, currentPath, "/etc/hostname")
	if result {
		t.Log("download success")
	}

	result, _ = UploadFile(sftpObj, path.Join(currentPath, "hostname"), "/home")
	if result {
		t.Log("upload success")
	}

}
