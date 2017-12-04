package main

import (
	"path"
	"github.com/cheneylew/goutil/utils"
	"strings"
	"fmt"
)

type MyReqeust struct {
	Url string
	RequestParams string
	Response string
}

func ReqeustToAPIMain()  {
	fPath := path.Join(utils.SelfDir(),"request.txt")
	requestStrs := strings.Split(utils.FileReadAllString(fPath), "===")
	requests := make([]MyReqeust, 0)
	for _, requestStr := range requestStrs {
		if len(requestStr) > 10 {
			strs := strings.Split(requestStr, "$$$")
			request := MyReqeust{
				Url:utils.Trim(strs[1]),
				RequestParams:utils.Trim(strs[2]),
				Response:utils.Trim(strs[3]),
			}
			requests = append(requests, request)
		}
	}

	text := ""
	for _, request := range requests {
		coms := strings.Split(request.Url,"/");
		method := coms[len(coms) - 1];
		method = "- (void)"+method+"With"
		kvs := strings.Split(request.RequestParams,"&")
		keys := []string{}
		for k, v := range kvs {
			param := strings.Split(v, "=")
			if len(param) == 2 {
				key := param[0]
				//value := param[1]
				if key != "token" && key != "appversion" && key != "netNo" {
					keys = append(keys, key)

					if k == 0 {
						method += strFirstToUpper(key)
					} else {
						method += key
					}

					method += ":(NSString *) "
					method += key
					method += "\n"
				}
			}
		}
		utils.JJKPrintln(request.Url)
		method += "completion:(void (^)(JJUser* dataModel))success \nfailure:(void (^)(NVErrorDataModel* dataModel))failure {\n"
		method += fmt.Sprintf("    NSString* path = @\"%s\";\n", strings.Replace(request.Url,"http://10.12.8.11:8080/api","", -1))
		method += `
    NSMutableDictionary *params = [NSMutableDictionary new];`
		for _, value := range keys {
			method += fmt.Sprintf("\n    [params nvSetObject:%s forKey:@\"%s\"];", value, value)
		}

		method += "\n"
		method += `    [self request:path
           mothed:@"POST"
          version:@"1.0"
       parameters:params
          success:^(id responseObject) {
              NSDictionary* dic = [responseObject valueForKey:kResponseData];

              if (success) {
                  JJUser * dataModel = [JJUser mj_objectWithKeyValues:dic];
                  success(dataModel);
              }
          }
          failure:^(NVErrorDataModel *dataModel) {
              if (failure) {
                  failure(dataModel);
              }
          }];
}`

		text += method
		text += "\n"
	}

	utils.FileWriteString(path.Join(utils.SelfDir(),"result.txt"), text)
}


