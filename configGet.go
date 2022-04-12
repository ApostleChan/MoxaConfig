package main

import (
	"log"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

type host struct {
	hostname string
	ip       string
}

func GetHostInfo() ([]host, error) {
	data := []host{}
	// 从配置文件中读取主机列表
	curPath, _ := os.Getwd()
	_, err := os.Stat(path.Join(curPath, "config.ini"))
	if err != nil {
		return []host{}, err
	}
	cfg, err := ini.Load("config.ini")
	getErr("load config", err)

	for _, v := range cfg.Sections() {
		if v.Key("hostname").String() != "" {
			data = append(data, host{v.Key("hostname").String(), v.Key("ip").String()})
		}
	}

	return data, nil

	// 遍历所有的section
	// for _, v := range cfg.Sections() {
	// 	fmt.Println(v.GetKey("hostname"))
	// 	hostname, _ := v.GetKey("hostname")
	// 	ip, _ := v.GetKey("ip")
	// 	fmt.Println(v.GetKey("ip"))
	// 	data = append(data, map[string]string{hostname: ip})
	// }

	// // 获取默认分区的key
	// fmt.Println(cfg.Section("").Key("version").String()) // 将结果转为string

	// // 获取mysql分区的key
	// fmt.Println(cfg.Section("Com").Key("host").String()) // 将结果转为string

	// // 如果读取的值不在候选列表内，则会回退使用提供的默认值
	// fmt.Println("Server Protocol:",
	// 	cfg.Section("mysql").Key("port").In("80", []string{"5555", "8080"}))

	// // 自动类型转换
	// fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("mysql").Key("port").MustInt(9999))
	// fmt.Printf("Database Name: (%[1]T) %[1]s\n", cfg.Section("mysql").Key("database").MustString("test"))

	// // 修改某个值然后进行保存
	// cfg.Section("").Key("version").SetValue("2.0.0")
	// cfg.SaveTo("config.ini")
	// time.Sleep(1000 * time.Second)
}

func getErr(msg string, err error) {
	if err != nil {
		log.Printf("%v err->%v\n", msg, err)
	}
}
