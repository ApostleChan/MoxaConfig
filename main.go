package main

import (
	"image/color"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 运行程序
// go run main.go configGet.go sshGet.go image.go font.go

// 打包程序
// fyne package -os windows -icon moxa.png

// 自定义字体，重写主题
type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

// return bundled font resource
// ResourceSourceHanSansTtf 即是 bundle.go 文件中 var 的变量名
func (m MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	return resourceAlibabaPuHuiTiMediumTtf
}
func (*MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

func main() {

	a := app.New()                    //创建应用程序
	a.Settings().SetTheme(&MyTheme{}) // 设置中文主题
	// NewWindow(title string) Window 创建一个新的窗口
	w := a.NewWindow("DA682B配置工具-v1.0.1")

	w.Resize(fyne.Size{400, 500})
	w.CenterOnScreen()
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
		widget.NewFormItem("主机名：", hostEntry),
		widget.NewFormItem("用户名：", userEntry),
		widget.NewFormItem("密码：", passwordEntry),
	)

	from.Items[0].HintText = "请输入主机名"
	from.Items[1].HintText = "请输入用户名"
	from.Items[2].HintText = "请输入密码"

	currentPath, _ := os.Getwd()

	//获取主机名配置
	getHostnameCfgBtn := widget.NewButton("获取主机名（/etc/hostname)配置", func() {
		hostname := hostEntry.Text
		user := userEntry.Text
		passwd := passwordEntry.Text
		remoteFilePath := "/etc/hostname"
		loaclDirName := hostEntry.Text
		go func() {
			stfpObj, err := NewSftpclient(user, passwd, hostname, 22)
			if err != nil {
				msgLabel.SetText("连接错误:" + err.Error() + "\n请检查连接信息是否正确！")
				return
			} else {
				_, err = DownnloadFile(stfpObj, path.Join(currentPath, loaclDirName), remoteFilePath)
				if err != nil {
					msgLabel.SetText("下载hostanme文件失败，远端文件路径:" + " " + remoteFilePath)
				} else {
					msgLabel.SetText("下载hostanme文件成功，本地文件路径:" + " " + path.Join(currentPath, loaclDirName, "hostname"))
					stfpObj.Close()
				}
			}
		}()
		return
	})

	// 获取网卡配置
	getNetworkCfgBtn := widget.NewButton("获取网络（/etc/network/interfaces)配置", func() {
		remoteFilePath := "/etc/network/interfaces"
		loaclDirName := hostEntry.Text
		hostname := hostEntry.Text
		user := userEntry.Text
		passwd := passwordEntry.Text
		go func() {
			stfpObj, err := NewSftpclient(user, passwd, hostname, 22)

			if err != nil {
				msgLabel.SetText("连接错误:" + err.Error() + "\n请检查连接信息是否正确！")
				return
			} else {
				_, err = DownnloadFile(stfpObj, path.Join(currentPath, loaclDirName), remoteFilePath)
				if err != nil {
					msgLabel.SetText("下载网络配置interfaces失败，远端文件路径:" + " " + remoteFilePath)
				} else {
					msgLabel.SetText("下载网络配置interfaces成功，本地文件路径:" + " " + path.Join(currentPath, loaclDirName, "interfaces"))
					stfpObj.Close()
				}
			}
		}()
	})

	// 更新主机配置
	updateHostnameCfgBtn := widget.NewButton("上传主机名（/etc/hostname)配置", func() {
		remotePath := "/etc"
		loaclDirName := hostEntry.Text
		fileName := "hostname"
		hostname := hostEntry.Text
		user := userEntry.Text
		passwd := passwordEntry.Text
		go func() {
			stfpObj, err := NewSftpclient(user, passwd, hostname, 22)

			if err != nil {
				msgLabel.SetText("连接错误:" + err.Error() + "\n请检查连接信息是否正确！")
				return
			} else {
				_, err = UploadFile(stfpObj, path.Join(currentPath, loaclDirName, fileName), remotePath)
				if err != nil {
					msgLabel.SetText("上传主机配置hostname失败，错误原因:" + " " + err.Error())
				} else {
					msgLabel.SetText("上传主机配置hostname成功，远端文件路径:" + " " + remotePath + fileName)
					stfpObj.Close()
				}
			}
		}()
	})
	// 更新网络配置
	updateNetworkCfgBtn := widget.NewButton("上传网络（/etc/network/interfaces)配置", func() {
		remotePath := "/etc/network"
		loaclDirName := hostEntry.Text
		fileName := "interfaces"
		hostname := hostEntry.Text
		user := userEntry.Text
		passwd := passwordEntry.Text
		go func() {
			stfpObj, err := NewSftpclient(user, passwd, hostname, 22)

			if err != nil {
				msgLabel.SetText("连接错误:" + err.Error() + "\n请检查连接信息是否正确！")
				return
			} else {
				_, err = UploadFile(stfpObj, path.Join(currentPath, loaclDirName, fileName), remotePath)
				if err != nil {
					msgLabel.SetText("上传网络配置interfaces失败，错误原因:" + " " + err.Error())
				} else {
					msgLabel.SetText("上传网络配置interfaces成功，远端文件路径:" + " " + remotePath + fileName)
					stfpObj.Close()
				}
			}
		}()
	})

	ExcuteShellCommandBtn := widget.NewButton("重启机器", func() {
		hostname := hostEntry.Text + ":22"
		user := userEntry.Text
		passwd := passwordEntry.Text
		cli := Cli{
			user: user,
			pwd:  passwd,
			addr: hostname,
		}
		cli.Run("reboot")
		msgLabel.SetText("重启命令已经执行！")
	})

	c := container.NewVBox(
		from,
		getHostnameCfgBtn,
		getNetworkCfgBtn,
		updateHostnameCfgBtn,
		updateNetworkCfgBtn,
		ExcuteShellCommandBtn,
		msgLabel,
	)
	w.SetContent(c)

	// // ShowAndRun is a shortcut to show the window and then run the application.
	w.ShowAndRun()
}
