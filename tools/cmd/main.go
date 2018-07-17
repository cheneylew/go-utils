package main

import (
	"github.com/cheneylew/goutil/utils"
	"fmt"
	"strings"
	"path"
)

func main() {
	msg := `[1]杀掉某个端口进程
[2]所有正在监听端口
[3]某个端口占用情况
[4]显示所有进程
[5]根据进程名称来kill`
	fmt.Println(msg)
	id := utils.InputIntWithMessage("请选择功能:")
	switch id {
	case 1:
		portNum :=utils.InputIntWithMessage("输入端口号:")
		if portNum > 0 {
			ports := getListeningPorts()
			for _, value := range ports {
				if value.PortNum == portNum {
					shell := fmt.Sprintf("kill %d", value.Pid)
					utils.ExecShell(shell)
					utils.JJKPrintln(fmt.Sprintf("%s success!", shell))
				}
			}
		}
	case 2:
		if utils.IsMac() {
			msg := utils.ExecShell(`lsof -i -P | grep -i "listen"`)
			utils.JJKPrintln(msg)
		} else if utils.IsLinux() {
			msg := utils.ExecShell(`netstat -ntlp`)
			utils.JJKPrintln(msg)
		}
	case 3:
		portNum :=utils.InputIntWithMessage("输入端口号:")
		if portNum > 0 {
			if utils.IsMac() {
				rst := utils.ExecShell(fmt.Sprintf(`lsof -i tcp:%d  | grep -i "listen"`, portNum))
				if len(rst) > 0 {
					utils.JJKPrintln(rst)
				}  else {
					utils.JJKPrintln(fmt.Sprintf("%d 此端口没有使用!", portNum))
				}
			} else if utils.IsLinux() {
				rst := utils.ExecShell(fmt.Sprintf(`netstat -tunlp|grep %d`, portNum))
				if len(rst) > 0 {
					utils.JJKPrintln(rst)
				}  else {
					utils.JJKPrintln(fmt.Sprintf("%d 此端口没有使用!", portNum))
				}
			}
		}
	case 4:
		utils.JJKPrintln(utils.ExecShell("ps -le"))
	case 5:
		pname := utils.InputStringWithMessage("输入进程名称：")
		rst := utils.ExecShell("ps -le")
		arr := strings.Split(rst, "\n")
		for index, row := range arr {
			if index > 0 {
				items := utils.RegexpFindAll(row,`\S+`)
				appPath := ""
				appName := ""
				pid := ""
				if utils.IsMac() {
					appPath = items[14]
					_, appName = path.Split(appPath)
					pid = items[1]
				} else if utils.IsLinux() {
					pid = items[3]
					appPath = items[len(items)-1]
					_, appName = path.Split(appPath)
				}

				if pname == appName {
					shell := fmt.Sprintf("kill %s", pid)
					utils.ExecShell(shell)
					utils.JJKPrintln(fmt.Sprintf("kill %s success!", pid))
				}
			}
		}
	}

	utils.JJKPrintln("finished!")
}

type Port struct {
	Pid int
	PortNum int
}

func getListeningPorts() []Port {
	var ports []Port
	if utils.IsMac() {
		rst := utils.ExecShell(`lsof -i -P | grep -i "listen"`)
		rows := strings.Split(rst,"\n")
		for _, row := range rows {
			arr := utils.RegexpFindAll(row,`\S+`)
			pid := arr[1]
			portNum := utils.LastOne(utils.RegexpFindAll(arr[8], `\w+`))
			utils.JJKPrintln(pid, portNum)
			ports = append(ports, Port{utils.JKStrToInt(pid),utils.JKStrToInt(portNum)})
		}
	} else if utils.IsLinux() {
		rst := utils.ExecShell(`netstat -ntlp`)
		rows := strings.Split(rst,"\n")
		for index, row := range rows {
			if index > 1 {
				arr := utils.RegexpFindAll(row, `\S+`)
				portStr := arr[3]
				pidStr := arr[len(arr)-1]
				portNum := utils.LastOne(utils.RegexpFindAll(portStr, `\w+`))
				pid := utils.RegexpFindAll(pidStr, `\w+`)[0]
				ports = append(ports, Port{utils.JKStrToInt(pid), utils.JKStrToInt(portNum)})
			}
		}
	}
	return ports
}