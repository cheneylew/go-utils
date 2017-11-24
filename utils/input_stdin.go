package utils

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
)

func InputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text,string('\n'),"",-1)
	return text
}

func InputInt() int  {
	n,_ := strconv.Atoi(InputString())
	return n;
}

func InputStringWithMessage(msg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text,string('\n'),"",-1)
	return text
}

func InputIntWithMessage(msg string) int {
	n,_ := strconv.Atoi(InputStringWithMessage(msg))
	return n;
}