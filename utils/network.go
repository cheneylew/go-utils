package utils

import (
	"fmt"
	"net"
	"strings"
)

func JKMACAddressEn0() string {
	result := ""
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return result
	}

	for _, inter := range interfaces {
		//fmt.Println(inter.Name)
		mac := inter.HardwareAddr //获取本机MAC地址
		macStr := fmt.Sprintf("%s", mac)
		name := strings.ToLower(inter.Name)
		if name == "en0" || strings.Contains(name, "本地连接") {
			result = macStr
			break
		}
	}

	return result
}

func JKGetLocalIP() string {
	en0Ip := JKGetIPWithInterface("en0")
	en1Ip := JKGetIPWithInterface("en1")
	if len(en0Ip) > 7 {
		return en0Ip
	}

	if len(en1Ip) > 7 {
		return en1Ip
	}

	return ""
}

func JKGetIPWithInterface(name string) string {
	ifi, _ := net.InterfaceByName(name)
	addrs, _ := ifi.Addrs()
	ip := ""
	for _, a := range addrs {
		ip = fmt.Sprintf("%s", a)
	}
	arr := strings.Split(ip,"/")
	if len(arr) >= 1 {
		return arr[0]
	}
	return ip
}


func JKIPPortIsListening(ip, port string)  bool {
	result := ExecShell("telnet "+ip+" "+port)
	if strings.Contains(result, "Connected") {
		return true
	} else {
		return false
	}
}
