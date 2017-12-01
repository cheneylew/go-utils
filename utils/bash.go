package utils

import (
	"os/exec"
	"strings"
	"fmt"
	"os"
)

func ExecComandWithDir(dir string, cmd string, params ...string) string {
	c := exec.Command(cmd, params...)
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		fmt.Printf("Exec error: %v; %v; %v; %v;\n", err.Error(),dir,cmd, params)
		return ""
	}
	s := string(out)
	return Trim(s)
}

func ExecComand(cmd string, params ...string) string {
	c := exec.Command(cmd, params...)
	out, _ := c.Output()
	s := string(out)
	s = strings.TrimSuffix(s, string('\n'))
	return s
}

func ExecShell(shell string) string {
	tmpShFile := ExeDirAppend("tmp.sh")
	FileWriteString(tmpShFile, fmt.Sprintf("#!/usr/bin/env bash \n%v", shell))
	result := ExecComand("sh","tmp.sh")
	os.Remove(tmpShFile)
	return result
}

func ExecShellInDir(dir, shell string) string {
	tmpShFile := ExeDirAppend("tmp.sh")
	FileWriteString(tmpShFile, fmt.Sprintf("#!/usr/bin/env bash \n%v", shell))
	result := ExecComandWithDir(dir,"sh","tmp.sh")
	os.Remove(tmpShFile)
	return result
}