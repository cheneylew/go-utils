package main

func main() {
	flag := 0
	if flag == 0 {
		//将后端BizConstant.java所有常量转换为iOS宏定义
		fileDir := "/Users/dejunliu/Desktop/"
		inputFileName := "BizConstant.java"
		outFileName := "CommonDefine.pch"
		JavaConstProcess(fileDir, inputFileName, outFileName)
	}
}
