package main

func main() {
	/*flag value
	//0.将后端BizConstant.java所有常量转换为iOS宏定义
	//1.将Json转为Objective-C的Model
	//2.将H5接口请求转为iOS接口请求
	 */
	flag := 2
	if flag == 0 {
		fileDir := "/Users/dejunliu/Desktop/"
		inputFileName := "BizConstant.java"
		outFileName := "CommonDefine.pch"
		JavaConstProcessMain(fileDir, inputFileName, outFileName)
	} else if flag == 1 {
		JsonToOCModelMain("/json.txt")
	} else if flag == 2 {
		ReqeustToAPIMain()
	}
}
