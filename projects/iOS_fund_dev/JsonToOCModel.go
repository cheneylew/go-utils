package main

import (
	"strings"
	"fmt"
	"github.com/cheneylew/goutil/utils"
	"encoding/json"
)

var hHeader string = `//
//  App
//
//  Created by Dejun Liu on 2017/11/24.
//  Copyright © 2017年 Dejun Liu. All rights reserved.
//

#import <Foundation/Foundation.h>

`

var mHeader string = `//
//
//  App
//
//  Created by Dejun Liu on 2017/11/24.
//  Copyright © 2017年 Dejun Liu. All rights reserved.
//

`

func JsonToOCModelMain(jsonFileName string) {
	jpath := utils.SelfDir()+jsonFileName
	jsonstr := utils.FileReadAllString(jpath)
	if len(jsonstr) == 0 {
		utils.JJKPrintln("请把需要转换的文本放到："+jpath)
		return ;
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonstr), &data)
	if err != nil {
		utils.JJKPrintln(err)
		return
	}

	innerData := data["data"]

	objectString = ""
	classNames = classNames[:0]

	name := utils.InputStringWithMessage("RootClassName:")
	if len(name) <= 2 {
		utils.JJKPrintln("类名太短了。。。")
		return
	}


	modelName := name
	className := getClassName(modelName)

	utils.JJKPrintln(className)

	str := recurse(innerData, modelName)
	if len(str) > 0 {

		fileText := hHeader + objectString
		//写.h
		utils.FileWriteString(fmt.Sprintf("%s/%s.h", utils.SelfDir(), className), fileText)

		impl := ""

		for _,v := range fileClasses {
			isContainPropertyArray := false
			for _, property := range v.Properties {
				if strings.Contains(property.ObjectCType, "NSArray") {
					isContainPropertyArray = true
				}
			}

			arrayProperties := ""
			if isContainPropertyArray {
				arrayProperties += `+ (NSDictionary *)mj_objectClassInArray {
    return @{`
				var strS []string
				for _, property := range v.Properties {
					if strings.Contains(property.ObjectCType, "NSArray") {
						strS = append(strS,fmt.Sprintf("@\"%s\":@\"%s\"", property.Name, property.ObjectArrayClass))
					}
				}

				arrayProperties += strings.Join(strS, ",\n")
				arrayProperties += `};
}`
			}

			if v.IsContainId {
				ids := `+ (NSDictionary *)mj_replacedKeyFromPropertyName {
    return @{@"ID":@"id"};
}`
				impl += fmt.Sprintf("@implementation %s \n%s\n%s\n@end\n",v.Name,ids,arrayProperties)
			} else {
				impl += fmt.Sprintf("@implementation %s \n%s\n@end\n", v.Name,arrayProperties)
			}

		}

		impl = mHeader +fmt.Sprintf("\n#import \"%s.h\"\n\n\n", className)+ impl

		//写.m
		utils.FileWriteString(fmt.Sprintf("%s/%s.m" , utils.SelfDir(),className), impl)
	}

	utils.ExecShell("open "+utils.SelfDir())
}


type Property struct {
	Type string
	ObjectCType string
	ObjectArrayClass string
	Name string
}

type Class struct {
	Name string
	IsContainId bool
	Properties []Property
}

var classPrefix string = "JJ"
var fileClasses []Class = make([]Class,0)
var objectString string = ""
var classNames []string


func strFirstToUpper(str string) string {
	if len(str) > 0 {
		if len(str) == 1 {
			return strings.ToUpper(str)
		} else {
			return strings.ToUpper(string(str[0]))+str[1:]
		}

	}
	return str
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func getClassName(name string) string {
	newName := strFirstToUpper(trimSuffix(name, "s"))
	return fmt.Sprintf("%s%s", classPrefix,newName)
}

func recurse(input interface{}, name string) string {
	switch input.(type) {
	case map[string]interface{}:
		str := ""
		className := getClassName(name)
		class := fmt.Sprintf("\n@interface %s : NSObject\n\n", className)
		classNames = append(classNames, className)
		mp, _ := input.(map[string]interface{})

		properties := make([]Property, 0)
		c := Class{
			Name:className,
			IsContainId:false,
			Properties:properties,
		}

		for key, value := range mp {
			if key == "id" {
				c.IsContainId = true
				key = "ID"
			}

			property := Property{
				Name:key,
				Type:fmt.Sprintf("%T", value),
			}

			str += fmt.Sprintf("%s:%T", key, value)
			switch value.(type) {
			case map[string]interface{}:
				property.ObjectCType = fmt.Sprintf("%s", getClassName(key))
				class += fmt.Sprintf("PP_STRONG(%s, %s)\n", getClassName(key), key)
				str += fmt.Sprintf("=<<%s>>", recurse(value, key))
			case []interface{}:
				property.ObjectCType = fmt.Sprintf("NSArray<%s *>", getClassName(key))
				property.ObjectArrayClass = getClassName(key);
				class += fmt.Sprintf("PP_STRONG(NSArray<%s *>, %s)\n", getClassName(key), key)
				value1, _ := value.([]interface{})
				if len(value1) > 0 {
					str += fmt.Sprintf("=<<%s>>", recurse(value1[0], key))
				}
			case float64:
				class += fmt.Sprintf("PP_ASSIGN_BASIC(double, %s)\n", key)
				property.ObjectCType = "double"
			case bool:
				class += fmt.Sprintf("PP_ASSIGN_BASIC(BOOL, %s)\n", key)
				property.ObjectCType = "BOOL"
			default:
				class += fmt.Sprintf("PP_STRONG(NSString, %s)\n", key)
				property.ObjectCType = "NSString"
				break
			}
			str += fmt.Sprintf("\n")

			properties = append(properties, property)
		}

		c.Properties = properties
		fileClasses = append(fileClasses,c)

		class += fmt.Sprintf("\n@end")
		objectString += class + "\n"

		return str
	case []interface{}:
		arr, _ := input.([]interface{})
		for _,value := range arr {
			switch value.(type) {
			case []interface{}:
				value1, _ := input.([]interface{})
				if len(value1) > 0 {
					return recurse(value1[0], name)
				}
			case map[string]interface{}:
				return recurse(value, name)
			default:
				break
			}
		}
	}

	return ""


}
