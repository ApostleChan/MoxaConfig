// fyne package -os windows -icon moxa.png
package main

import (
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	a := app.New() //创建应用程序

	// NewWindow(title string) Window 创建一个新的窗口
	w := a.NewWindow("Moxa Config-v1.0.0")

	w.Resize(fyne.Size{400, 500})
	w.SetFixedSize(false)
	w.SetIcon(resourceMoxaPng)

	hostinfo, err := GetHostInfo()
	if err != nil {

	}
	ip_list := []string{}
	for _, v := range hostinfo {
		ip_list = append(ip_list, v.ip)
	}

	hostEntry := widget.NewSelectEntry(ip_list) //主机名
	userEntry := widget.NewEntry()              //用户名
	passwordEntry := widget.NewPasswordEntry()  //密码

	msgLabel := widget.NewLabel("")
	msgLabel.Wrapping = fyne.TextWrap(fyne.TextWrapWord)

	from := widget.NewForm(
		widget.NewFormItem("hostname:", hostEntry),
		widget.NewFormItem("user:", userEntry),
		widget.NewFormItem("passowrd:", passwordEntry),
	)

	from.Items[0].HintText = "please input hostname"
	from.Items[1].HintText = "please input username"
	from.Items[2].HintText = "please input password"

	currentPath, _ := os.Getwd()

	//获取主机名配置
	getHostnameCfgBtn := widget.NewButton("Get Hostname Config", func() {
		hostname := hostEntry.Text
		user := userEntry.Text
		passwd := passwordEntry.Text
		remoteFilePath := "/etc/hostname"
		loaclDirName := hostEntry.Text
		go func() {
			stfpObj, err := NewSftpclient(user, passwd, hostname, 22)
			if err != nil {
				msgLabel.SetText("connect error :" + err.Error())
				return
			} else {
				_, err = DownnloadFile(stfpObj, path.Join(currentPath, loaclDirName), remoteFilePath)
				if err != nil {
					msgLabel.SetText("download hostname config failed:" + " " + remoteFilePath)
				} else {
					msgLabel.SetText("download hostname config successed:" + " " + remoteFilePath)
					stfpObj.Close()
				}
			}
		}()
		return
	})

	// 获取网卡配置
	getNetworkCfgBtn := widget.NewButton("Get Network Config", func() {
		remoteFilePath := "/etc/network/interfaces"
		loaclDirName := hostEntry.Text
		hostname := hostEntry.Text
		user := userEntry.Text
		passwd := passwordEntry.Text
		go func() {
			stfpObj, err := NewSftpclient(user, passwd, hostname, 22)

			if err != nil {
				msgLabel.SetText("connect error :" + err.Error())
				return
			} else {
				_, err = DownnloadFile(stfpObj, path.Join(currentPath, loaclDirName), remoteFilePath)
				if err != nil {
					msgLabel.SetText("download network config failed:" + " " + remoteFilePath)
				} else {
					msgLabel.SetText("download network config successed:" + " " + remoteFilePath)
					stfpObj.Close()
				}
			}
		}()
	})

	// 更新主机配置
	updateHostnameCfgBtn := widget.NewButton("Update Hostname Config", func() {
		remotePath := "/etc"
		loaclDirName := hostEntry.Text
		fileName := "hostname1"
		hostname := hostEntry.Text
		user := userEntry.Text
		passwd := passwordEntry.Text
		go func() {
			stfpObj, err := NewSftpclient(user, passwd, hostname, 22)

			if err != nil {
				msgLabel.SetText("connect error :" + err.Error())
				return
			} else {
				_, err = UploadFile(stfpObj, path.Join(currentPath, loaclDirName, fileName), remotePath)
				if err != nil {
					msgLabel.SetText("upload hostname config failed:" + " " + err.Error())
				} else {
					msgLabel.SetText("upload hostname config successed:" + " " + remotePath)
					stfpObj.Close()
				}
			}
		}()
	})
	// 更新网络配置
	updateNetworkCfgBtn := widget.NewButton("Update Networking Config", func() {
		remotePath := "/etc/network"
		loaclDirName := hostEntry.Text
		fileName := "interfaces1"
		hostname := hostEntry.Text
		user := userEntry.Text
		passwd := passwordEntry.Text
		go func() {
			stfpObj, err := NewSftpclient(user, passwd, hostname, 22)

			if err != nil {
				msgLabel.SetText("connect error :" + err.Error())
				return
			} else {
				_, err = UploadFile(stfpObj, path.Join(currentPath, loaclDirName, fileName), remotePath)
				if err != nil {
					msgLabel.SetText("upload network config failed:" + " " + err.Error())
				} else {
					msgLabel.SetText("upload network config successed:" + " " + remotePath)
					stfpObj.Close()
				}
			}
		}()
	})

	c := container.NewVBox(from, getHostnameCfgBtn, getNetworkCfgBtn, updateHostnameCfgBtn, updateNetworkCfgBtn, msgLabel)
	w.SetContent(c)

	// // ShowAndRun is a shortcut to show the window and then run the application.
	w.ShowAndRun()
}
